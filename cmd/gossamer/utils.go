// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/ChainSafe/gossamer/dot"
	ctoml "github.com/ChainSafe/gossamer/dot/config/toml"
	"github.com/ChainSafe/gossamer/internal/log"
	"github.com/urfave/cli"
	terminal "golang.org/x/term"
)

const confirmCharacter = "Y"

// setupLogger sets up the global Gossamer logger.
func setupLogger(ctx *cli.Context) (level log.Level, err error) {
	level, err = getLogLevel(ctx, LogFlag.Name, "", log.Info)
	if err != nil {
		return level, err
	}

	log.Patch(
		log.SetWriter(os.Stdout),
		log.SetFormat(log.FormatConsole),
		log.SetCallerFile(true),
		log.SetCallerLine(true),
		log.SetLevel(level),
	)

	return level, nil
}

// getPassword prompts user to enter password
func getPassword(msg string) []byte {
	for {
		fmt.Println(msg)
		fmt.Print("> ")
		password, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Printf("invalid input: %s\n", err)
		} else {
			fmt.Printf("\n")
			return password
		}
	}
}

// confirmMessage prompts user to confirm message and returns true if "Y"
func confirmMessage(msg string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(msg)
	fmt.Print("> ")
	for {
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		return strings.Compare(confirmCharacter, strings.ToUpper(text)) == 0
	}
}

// newTestConfig returns a new test configuration using the provided basepath
func newTestConfig(t *testing.T) *dot.Config {
	dir := t.TempDir()

	cfg := &dot.Config{
		Global: dot.GlobalConfig{
			Name:           dot.GssmrConfig().Global.Name,
			ID:             dot.GssmrConfig().Global.ID,
			BasePath:       dir,
			LogLvl:         log.Info,
			PublishMetrics: dot.GssmrConfig().Global.PublishMetrics,
			MetricsPort:    dot.GssmrConfig().Global.MetricsPort,
			RetainBlocks:   dot.GssmrConfig().Global.RetainBlocks,
			Pruning:        dot.GssmrConfig().Global.Pruning,
			TelemetryURLs:  dot.GssmrConfig().Global.TelemetryURLs,
		},
		Log: dot.LogConfig{
			CoreLvl:           log.Info,
			DigestLvl:         log.Info,
			SyncLvl:           log.Info,
			NetworkLvl:        log.Info,
			RPCLvl:            log.Info,
			StateLvl:          log.Info,
			RuntimeLvl:        log.Info,
			BlockProducerLvl:  log.Info,
			FinalityGadgetLvl: log.Info,
		},
		Init:    dot.GssmrConfig().Init,
		Account: dot.GssmrConfig().Account,
		Core:    dot.GssmrConfig().Core,
		Network: dot.GssmrConfig().Network,
		RPC:     dot.GssmrConfig().RPC,
		System:  dot.GssmrConfig().System,
		Pprof:   dot.GssmrConfig().Pprof,
	}

	return cfg
}

// newTestConfigWithFile returns a new test configuration and a temporary configuration file
func newTestConfigWithFile(t *testing.T) (*dot.Config, *os.File) {
	cfg := newTestConfig(t)

	filename := filepath.Join(cfg.Global.BasePath, "config.toml")

	tomlCfg := dotConfigToToml(cfg)
	cfgFile := exportConfig(tomlCfg, filename)
	return cfg, cfgFile
}

func dotConfigToToml(dcfg *dot.Config) *ctoml.Config {
	cfg := &ctoml.Config{
		Pprof: ctoml.PprofConfig{
			Enabled:          dcfg.Pprof.Enabled,
			ListeningAddress: dcfg.Pprof.Settings.ListeningAddress,
			BlockRate:        dcfg.Pprof.Settings.BlockProfileRate,
			MutexRate:        dcfg.Pprof.Settings.MutexProfileRate,
		},
	}

	cfg.Global = ctoml.GlobalConfig{
		Name:         dcfg.Global.Name,
		ID:           dcfg.Global.ID,
		BasePath:     dcfg.Global.BasePath,
		LogLvl:       dcfg.Global.LogLvl.String(),
		MetricsPort:  dcfg.Global.MetricsPort,
		RetainBlocks: dcfg.Global.RetainBlocks,
		Pruning:      string(dcfg.Global.Pruning),
	}

	cfg.Log = ctoml.LogConfig{
		CoreLvl:           dcfg.Log.CoreLvl.String(),
		SyncLvl:           dcfg.Log.SyncLvl.String(),
		NetworkLvl:        dcfg.Log.NetworkLvl.String(),
		RPCLvl:            dcfg.Log.RPCLvl.String(),
		StateLvl:          dcfg.Log.StateLvl.String(),
		RuntimeLvl:        dcfg.Log.RuntimeLvl.String(),
		BlockProducerLvl:  dcfg.Log.BlockProducerLvl.String(),
		FinalityGadgetLvl: dcfg.Log.FinalityGadgetLvl.String(),
	}

	cfg.Init = ctoml.InitConfig{
		Genesis: dcfg.Init.Genesis,
	}

	cfg.Account = ctoml.AccountConfig{
		Key:    dcfg.Account.Key,
		Unlock: dcfg.Account.Unlock,
	}

	cfg.Core = ctoml.CoreConfig{
		Roles:            dcfg.Core.Roles,
		BabeAuthority:    dcfg.Core.BabeAuthority,
		GrandpaAuthority: dcfg.Core.GrandpaAuthority,
		GrandpaInterval:  uint32(dcfg.Core.GrandpaInterval / time.Second),
	}

	cfg.Network = ctoml.NetworkConfig{
		Port:              dcfg.Network.Port,
		Bootnodes:         dcfg.Network.Bootnodes,
		ProtocolID:        dcfg.Network.ProtocolID,
		NoBootstrap:       dcfg.Network.NoBootstrap,
		NoMDNS:            dcfg.Network.NoMDNS,
		DiscoveryInterval: int(dcfg.Network.DiscoveryInterval / time.Second),
		MinPeers:          dcfg.Network.MinPeers,
		MaxPeers:          dcfg.Network.MaxPeers,
	}

	cfg.RPC = ctoml.RPCConfig{
		Enabled:    dcfg.RPC.Enabled,
		External:   dcfg.RPC.External,
		Port:       dcfg.RPC.Port,
		Host:       dcfg.RPC.Host,
		Modules:    dcfg.RPC.Modules,
		WSPort:     dcfg.RPC.WSPort,
		WS:         dcfg.RPC.WS,
		WSExternal: dcfg.RPC.WSExternal,
	}

	return cfg
}
