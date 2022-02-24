package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"github.com/viant/toolbox"
	"gopkg.in/yaml.v2"
	"github.com/viant/smirror/config"
	"github.com/viant/smirror/shared"
	"time"
)

func (s *service) loadRule(ctx context.Context, URL string) (*config.Rule, error) {
	reader, err := s.fs.OpenURL(ctx, URL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download rule: %v", URL)
	}
	defer reader.Close()
	_, name := url.Split(URL, "")
	ruleURL := url.Join(s.config.Mirrors.BaseURL, name)
	err = s.fs.Upload(ctx, ruleURL, file.DefaultFileOsMode, reader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update rule: %v", ruleURL)
	}

	//add sleep to refresh
	time.Sleep(1 * time.Millisecond)
	_, err = s.config.Mirrors.ReloadIfNeeded(ctx, s.fs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to reload rules: %v", ruleURL)
	}
	rule := s.config.Mirrors.Rule(ruleURL)
	if rule == nil {
		return nil, errors.Errorf("failed to lookup rule: %v", ruleURL)
	}
	//Assume that rule is never disable for stand alone process
	rule.Disabled = false
	return rule, nil
}

func ruleToMap(rule *config.Rule) map[string]interface{} {
	ruleMap := map[string]interface{}{}
	toolbox.DefaultConverter.AssignConverted(&ruleMap, rule)
	compactedMap := map[string]interface{}{}
	toolbox.CopyMap(ruleMap, compactedMap, toolbox.OmitEmptyMapWriter)
	return compactedMap
}

func (s *service) reportRule(rule *config.Rule) {
	ruleMap := ruleToMap(rule)
	var ruleYAML, err = yaml.Marshal(ruleMap)
	if err == nil {
		shared.LogF("==== USING RULE ===\n%s===== END ====\n", ruleYAML)
	}
}

