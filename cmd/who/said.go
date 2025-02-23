package who

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"sync"

	"github.com/jxsl13/cli-config-boilerplate/cliconfig"
	"github.com/jxsl13/twlog/cmd/who/config"
	"github.com/jxsl13/twlog/ctxutils"
	"github.com/jxsl13/twlog/fswalk"
	"github.com/jxsl13/twlog/internal/sharedcontext"
	"github.com/jxsl13/twlog/match"
	"github.com/jxsl13/twlog/model"
	"github.com/spf13/cobra"
)

func NewSaidCommand(root *sharedcontext.Root) *cobra.Command {
	cli := &SaidContext{
		root: root,
		cfg:  config.NewSaidConfig(),
	}

	cmd := cobra.Command{
		Use:   "said",
		Short: "said searches for what players said in the chat",
	}
	cmd.PreRunE = cli.PreRunE(&cmd)
	cmd.RunE = cli.RunE
	return &cmd
}

type SaidContext struct {
	root               *sharedcontext.Root
	cfg                config.SaidConfig
	SearchPhraseRegexp *regexp.Regexp
}

func (cli *SaidContext) PreRunE(cmd *cobra.Command) func(*cobra.Command, []string) error {
	parser := cliconfig.RegisterFlags(&cli.cfg, false, cmd, cliconfig.WithoutConfigFile())
	return func(cmd *cobra.Command, args []string) error {
		log.SetOutput(cmd.ErrOrStderr()) // redirect log output to stderr

		if len(args) == 0 {
			return errors.New("missing search phrase regex argument")
		}

		phrase, err := regexp.Compile(args[0])
		if err != nil {
			return fmt.Errorf("could not compile search phrase regex: %w", err)
		}

		cli.SearchPhraseRegexp = phrase

		return parser() // parse registered commands
	}
}

func (cli *SaidContext) RunE(cmd *cobra.Command, args []string) error {

	var (
		ctx                = cli.root.Ctx
		mu                 = &sync.Mutex{}
		extendedPlayerList = make(model.PlayerExtendedList, 0, 64)
		format             = cli.root.Format
	)

	err := fswalk.Walk(ctx, cli.root.Walk.ToFSWalkConfig(), func(filePath string, file io.Reader) error {
		filePlayers, err := searchPhrase(filePath, file, cli.SearchPhraseRegexp)
		if err != nil {

			return err
		}
		mu.Lock()
		defer mu.Unlock()
		extendedPlayerList = append(extendedPlayerList, filePlayers...)
		return nil
	})
	if err != nil {
		return err
	}

	err = ctxutils.Done(ctx)
	if err != nil {
		return err
	}

	if cli.cfg.IPsOnly {
		ipList := extendedPlayerList.ToIPList()
		if cli.cfg.Deduplicate {
			ipList = deduplicate(ipList)
		}
		return format.Print(cmd, ipList)
	} else if cli.cfg.Extended {
		if cli.cfg.Deduplicate {
			extendedPlayerList = deduplicate(extendedPlayerList)
		}
		return format.Print(cmd, extendedPlayerList)
	}

	// not extended list of players
	playerList := extendedPlayerList.ToPlayerList()
	if cli.cfg.Deduplicate {
		playerList = deduplicate(playerList)
	}

	return format.Print(cmd, playerList)
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

func searchPhrase(filePath string, f io.Reader, phraseRegexp *regexp.Regexp) (model.PlayerExtendedList, error) {

	players := make(model.PlayerExtendedList, 0, 16)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// id -> ip
	playerMap := make(map[int]string, 64)

	for scanner.Scan() {
		line := scanner.Text()

		if id, ip, ok := match.Join(line); ok {
			playerMap[id] = ip
			continue
		} else if id, ok := match.Leave(line); ok {
			delete(playerMap, id)
			continue
		} else if id, nick, chat, ok := match.Chat(line); ok {
			if !phraseRegexp.MatchString(chat) {
				continue
			}

			ip, ok := playerMap[id]
			if !ok {
				fmt.Printf("could not find join line for player %s with id: %d\n", nick, id)
				continue
			}

			players = append(players, model.PlayerExtended{
				File:     filePath,
				Nickname: nick,
				ID:       id,
				IP:       ip,
				Text:     chat,
			})

		}
	}

	if err := scanner.Err(); err != nil {
		if !errors.Is(err, io.EOF) {
			return players, err
		}
	}

	return players, nil
}
