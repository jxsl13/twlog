package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/icza/backscanner"
	"github.com/jxsl13/cli-config-boilerplate/cliconfig"
	"github.com/jxsl13/twlog-who-said/config"
	"github.com/spf13/cobra"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cmd := NewRootCmd(ctx)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func NewRootCmd(ctx context.Context) *cobra.Command {
	cli := &CLI{
		ctx: ctx,
		cfg: config.NewConfig(),
	}

	cmd := cobra.Command{
		Use: filepath.Base(os.Args[0]),
	}
	cmd.PreRunE = cli.PrerunE(&cmd)
	cmd.RunE = cli.RunE
	return &cmd
}

type CLI struct {
	ctx context.Context
	cfg config.Config
}

func (cli *CLI) PrerunE(cmd *cobra.Command) func(*cobra.Command, []string) error {
	parser := cliconfig.RegisterFlags(&cli.cfg, false, cmd)
	return func(cmd *cobra.Command, args []string) error {
		log.SetOutput(cmd.OutOrStdout()) // redirect log output to stdout
		return parser()                  // parse registered commands
	}
}

func (cli *CLI) RunE(cmd *cobra.Command, args []string) error {
	files := make([]string, 0, 16)

	entryDir := cli.cfg.SearchDir
	entryDir, err := filepath.Abs(entryDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of search dir: %w", err)
	}

	err = filepath.WalkDir(entryDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip non-files
		if !info.Type().IsRegular() {
			return nil
		}

		if !cli.cfg.FileRegexp.MatchString(path) {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}
	slices.Sort(files)

	extendedPlayerList := make(PlayerExtendedList, 0, 16)

	for _, file := range files {
		filePlayers, err := searchPhrase(file, cli.cfg.PhraseRegexp)
		if err != nil {
			return fmt.Errorf("failed to search phrase in file %s: %w", file, err)
		}
		extendedPlayerList = append(extendedPlayerList, filePlayers...)
	}

	if cli.cfg.IPsOnly {
		ipList := extendedPlayerList.ToIPList()
		if cli.cfg.Deduplicate {
			ipList = deduplicate(ipList)
		}
		return cli.print(cmd, ipList)
	} else if cli.cfg.Extended {
		if cli.cfg.Deduplicate {
			extendedPlayerList = deduplicate(extendedPlayerList)
		}
		return cli.print(cmd, extendedPlayerList)
	}

	// not extended list of players
	playerList := extendedPlayerList.ToPlayerList()
	if cli.cfg.Deduplicate {
		playerList = deduplicate(playerList)
	}

	return cli.print(cmd, playerList)
}

func (cli *CLI) print(cmd *cobra.Command, a any) error {
	switch cli.cfg.Output {
	case config.FormatText:
		return cli.printText(cmd, a)
	case config.FormatJSON:
		return cli.printJSON(cmd, a)
	default:
		// should never happen
		return fmt.Errorf("unsupported output format: %s", cli.cfg.Output)
	}
}

func (cli *CLI) printText(cmd *cobra.Command, a any) error {
	s := a.(fmt.Stringer) // will panic if used incorrectly
	_, err := fmt.Fprintln(cmd.OutOrStdout(), s.String())
	return err
}

func (cli *CLI) printJSON(cmd *cobra.Command, a any) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json result: %w", err)
	}

	_, err = cmd.OutOrStdout().Write(data)
	if err != nil {
		return fmt.Errorf("failed to print json result: %w", err)
	}
	fmt.Fprint(cmd.OutOrStdout(), "\n")
	return nil
}

func deduplicate[C comparable](items []C) []C {
	seen := make(map[C]struct{}, max(16, len(items)/16))
	unique := make([]C, 0, len(items))

	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		unique = append(unique, item)
	}
	return unique
}

var (
	// id, nick, chat line
	chatLineRegexp = regexp.MustCompile(`chat: (\d+):-?\d+:(.+): (.+)`)
)

func searchPhrase(file string, phraseRegexp *regexp.Regexp) (PlayerExtendedList, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	players := make(PlayerExtendedList, 0, 16)

	beginSearchOffset := 0
	scanner := bufio.NewScanner(f)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		beginSearchOffset = -len(data)
		return advance, token, err
	})

	for scanner.Scan() {
		line := scanner.Text()
		matches := chatLineRegexp.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}

		chat := matches[3]
		if !phraseRegexp.MatchString(chat) {
			continue
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			// must match, otherwise hte regex is wrong
			panic(err)
		}

		nick := matches[2]

		offset, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return players, fmt.Errorf("failed to get current offset: %w", err)
		}

		ip, ok, err := seekJoinLineBackwards(f, offset, beginSearchOffset, id)
		if err != nil {
			return players, err
		}

		if !ok {
			fmt.Printf("could not find join line for player %s with id: %d\n", nick, id)
			continue
		}

		players = append(players, PlayerExtended{
			File:     file,
			Nickname: nick,
			ID:       id,
			IP:       ip,
			Text:     chat,
		})
	}

	if err := scanner.Err(); err != nil {
		if !errors.Is(err, io.EOF) {
			return players, err
		}
	}

	return players, nil
}

func seekJoinLineBackwards(f *os.File, resetOffset int64, beginSearchOffset int, id int) (ip string, ok bool, err error) {
	defer func() {
		// return back to the position from which we started searching backwards
		_, returnErr := f.Seek(resetOffset, io.SeekStart)
		if returnErr != nil {
			err = errors.Join(err, returnErr)
		}
	}()

	// begin searching in reverse before the matched line
	scanOffset := int(resetOffset + int64(beginSearchOffset))
	backScanner := backscanner.New(f, scanOffset)
	for {
		line, _, err := backScanner.Line()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", false, nil
			}
			return "", false, err
		}

		ip, ok := matchJoinLineWithID(line, id)
		if ok {
			return ip, true, nil
		}
	}
}

func matchJoinLineWithID(line string, id int) (ip string, ok bool) {

	var (
		joinIDStr string
		joinIP    string
	)
	if matches := ddnetJoinRegex.FindStringSubmatch(line); len(matches) != 0 {
		joinIDStr = matches[1]
		joinIP = matches[2]
	} else if matches := playerzCatchJoinRegex.FindStringSubmatch(line); len(matches) != 0 {
		joinIDStr = matches[1]
		joinIP = matches[2]
	} else if matches := playerVanillaJoinRegex.FindStringSubmatch(line); len(matches) != 0 {
		joinIDStr = matches[1]
		joinIP = matches[2]
	} else {
		return "", false
	}

	joinId, err := strconv.Atoi(joinIDStr)
	if err != nil {
		return "", false
	}
	if joinId != id {
		return "", false
	}
	return joinIP, true
}

var (
	// 0: full 1: ID 2: IP
	ddnetJoinRegex = regexp.MustCompile(`(?i)player has entered the game\. ClientID=([\d]+) addr=[^\d]{0,2}([\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3})[^\d]{0,2}`)

	// 0: full 1: ID 2: IP 3: port 4: version 5: name 6: clan 7: country
	playerzCatchJoinRegex = regexp.MustCompile(`(?i)id=([\d]+) addr=([a-fA-F0-9\.\:\[\]]+):([\d]+) version=(\d+) name='(.{0,20})' clan='(.{0,16})' country=([-\d]+)$`)

	// 0: full 1: ID 2: IP
	playerVanillaJoinRegex = regexp.MustCompile(`(?i)player is ready\. ClientID=([\d]+) addr=[^\d]{0,2}([\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3})[^\d]{0,2}`)
)

type PlayerExtended struct {
	File     string `json:"file"`
	Nickname string `json:"nickname"`
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	Text     string `json:"text"`
}

func (p PlayerExtended) String() string {
	return fmt.Sprintf("%s:<%s> {%d} %s: %s", p.File, p.IP, p.ID, p.Nickname, p.Text)
}

type PlayerExtendedList []PlayerExtended

func (p PlayerExtendedList) String() string {
	var sb strings.Builder
	sb.Grow(len(p) * 512)
	for _, player := range p {
		sb.WriteString(player.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (p PlayerExtendedList) ToPlayerList() PlayerList {
	players := make([]Player, 0, len(p))
	for _, player := range p {
		players = append(players, Player{
			Nickname: player.Nickname,
			IP:       player.IP,
			Text:     player.Text,
		})
	}
	return players
}

func (p PlayerExtendedList) ToIPList() StringList {
	ips := make(StringList, 0, len(p))
	for _, player := range p {
		ips = append(ips, player.IP)
	}
	return ips
}

type Player struct {
	Nickname string `json:"nickname"`
	IP       string `json:"ip"`
	Text     string `json:"text"`
}

func (p Player) String() string {
	return fmt.Sprintf("<{%s}> %s: %s", p.IP, p.Nickname, p.Text)
}

type PlayerList []Player

func (p PlayerList) String() string {
	var sb strings.Builder
	sb.Grow(len(p) * 256)
	for _, player := range p {
		sb.WriteString(player.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}

type StringList []string

func (s StringList) String() string {
	var sb strings.Builder
	sb.Grow(len(s) * 64)
	for _, str := range s {
		sb.WriteString(str)
		sb.WriteByte('\n')
	}
	return sb.String()
}
