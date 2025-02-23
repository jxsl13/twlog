package what

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"sync"

	"github.com/jxsl13/cli-config-boilerplate/cliconfig"
	"github.com/jxsl13/twlog/config"
	"github.com/jxsl13/twlog/ctxutils"
	"github.com/jxsl13/twlog/fswalk"
	"github.com/jxsl13/twlog/internal/sharedcontext"
	"github.com/jxsl13/twlog/internal/sliceutils"
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
		Use:   "said [nickname regex]",
		Short: "said searches for what players said in the chat",
	}
	cmd.PreRunE = cli.PreRunE(&cmd)
	cmd.RunE = cli.RunE
	return &cmd
}

type SaidContext struct {
	root                 *sharedcontext.Root
	cfg                  config.SaidConfig
	NicknameSearchPhrase *regexp.Regexp
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
			return fmt.Errorf("could not compile nickname search phrase regex: %w", err)
		}

		cli.NicknameSearchPhrase = phrase

		return parser()
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
		filePlayers, err := searchNicknamePhrase(ctx, filePath, file, cli.NicknameSearchPhrase)
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
		ipTextList := extendedPlayerList.ToIPTextList()
		if cli.cfg.Deduplicate {
			ipTextList = sliceutils.Deduplicate(ipTextList)
		}
		return format.Print(cmd, ipTextList)
	} else if cli.cfg.Extended {
		if cli.cfg.Deduplicate {
			extendedPlayerList = sliceutils.Deduplicate(extendedPlayerList)
		}
		return format.Print(cmd, extendedPlayerList)
	}

	// not extended list of players
	playerList := extendedPlayerList.ToPlayerList()
	if cli.cfg.Deduplicate {
		playerList = sliceutils.Deduplicate(playerList)
	}

	return format.Print(cmd, playerList)
}

func searchNicknamePhrase(ctx context.Context, filePath string, f io.Reader, nicknameRegexp *regexp.Regexp) (model.PlayerExtendedList, error) {

	players := make(model.PlayerExtendedList, 0, 16)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// id -> ip
	playerMap := make(map[int]string, 64)
	var err error
	for scanner.Scan() {
		err = ctxutils.Done(ctx)
		if err != nil {
			return players, err
		}

		line := scanner.Text()

		if id, ip, ok := match.Join(line); ok {
			playerMap[id] = ip
			continue
		} else if id, ok := match.Leave(line); ok {
			delete(playerMap, id)
			continue
		} else if id, nick, chat, ok := match.Chat(line); ok {
			if !nicknameRegexp.MatchString(nick) {
				continue
			}

			ip, ok := playerMap[id]
			if !ok {
				fmt.Printf("could not find join line for player %s with id: %d\n", nick, id)
				continue
			}

			players = append(players, model.NewPlayerExtended(filePath, nick, id, ip, chat))
		}
	}

	if err := scanner.Err(); err != nil {
		if !errors.Is(err, io.EOF) {
			return players, err
		}
	}

	return players, nil
}
