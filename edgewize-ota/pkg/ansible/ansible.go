package ansible

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/apenella/go-ansible/v2/pkg/execute"
	"github.com/apenella/go-ansible/v2/pkg/execute/measure"
	results "github.com/apenella/go-ansible/v2/pkg/execute/result/json"
	"github.com/apenella/go-ansible/v2/pkg/execute/stdoutcallback"
	"github.com/apenella/go-ansible/v2/pkg/playbook"

	"k8s.io/klog/v2"
)

// Run runs an ansible playbook
func Run(connection, inventory, playbookFile string) (*string, error) {

	var err error
	var res *results.AnsiblePlaybookJSONResults

	buff := new(bytes.Buffer)

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory:  inventory,
		Connection: connection,
	}

	playbookCmd := playbook.NewAnsiblePlaybookCmd(playbook.WithPlaybooks(playbookFile), playbook.WithPlaybookOptions(ansiblePlaybookOptions))

	exec := measure.NewExecutorTimeMeasurement(
		stdoutcallback.NewJSONStdoutCallbackExecute(
			execute.NewDefaultExecute(
				execute.WithCmd(playbookCmd),
				execute.WithErrorEnrich(playbook.NewAnsiblePlaybookErrorEnrich()),
				execute.WithWrite(io.Writer(buff)),
			),
		),
	)

	err = exec.Execute(context.TODO())
	if err != nil {
		klog.Errorf("Failed to execute ansible playbook: %v", err)
	}

	klog.Infof("buff: %v", buff.String())

	res, err = results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		klog.Errorf("Failed to parse ansible playbook result: %v", err)
		return nil, err

	}
	return resString(res)
}

// 输出待优化
func resString(r *results.AnsiblePlaybookJSONResults) (*string, error) {
	str := ""
	for _, play := range r.Plays {
		for _, task := range play.Tasks {
			name := task.Task.Name
			for host, result := range task.Hosts {
				resultStr := ItemString(result)
				str = fmt.Sprintf("%s[%s] (%s)	%s\n", str, host, name, resultStr)
			}
		}
	}

	for host, stats := range r.Stats {
		if stats.Failures > 0 {
			str = fmt.Sprintf("%s[%s] failed", str, host)
			return &str, fmt.Errorf("host %s failed", host)
		}
		if stats.Unreachable > 0 {
			str = fmt.Sprintf("%s[%s] unreachable", str, host)
			return &str, fmt.Errorf("host %s unreachable", host)
		}
	}
	return &str, nil
}

func ItemString(res *results.AnsiblePlaybookJSONResultsPlayTaskHostsItem) string {
	return fmt.Sprintf("stdout: %v, stderr: %v", res.Stdout, res.Stderr)
}
