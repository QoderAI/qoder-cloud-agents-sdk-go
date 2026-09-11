package live

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

const ProjectMemoryPath = "projects/release-conventions.md"

// ProjectMemory keeps the remembered facts separate from the question. Varying
// the facts on each run prevents a plausible stock answer from passing.
type ProjectMemory struct {
	Project, ReleaseTime, Contact, RollbackVersion string
}

func NewProjectMemory() ProjectMemory {
	contacts := []string{"林岚", "陈朔", "叶澄", "苏棠"}
	return ProjectMemory{
		Project:         "青禾订单-" + Marker()[:6],
		ReleaseTime:     fmt.Sprintf("%02d:%02d", 20+rand.IntN(4), rand.IntN(60)),
		Contact:         contacts[rand.IntN(len(contacts))],
		RollbackVersion: fmt.Sprintf("v2.%d.%d", 100+rand.IntN(900), 100+rand.IntN(900)),
	}
}

func (m ProjectMemory) Content() string {
	return fmt.Sprintf("---\nname: release-conventions\ndescription: %s 的项目发布约定\nmetadata:\n  type: project\n---\n\n# %s 的发布约定\n\n- 这是一个订单服务项目。\n- 团队约定在北京时间 %s 开始发布。\n- 发布异常时先联系值班负责人%s。\n- 如果需要回滚，使用已验证的稳定版本 %s。\n\n**Why:** 团队需要在值班人员在岗的窗口发布，并使用验证过的版本恢复服务。\n**How to apply:** 为这个项目拟定发布计划时，遵循以上团队约定。\n", m.Project, m.Project, m.ReleaseTime, m.Contact, m.RollbackVersion)
}

// The runtime loads MEMORY.md as an index. Keep the facts in the linked entry
// so the Agent must actually retrieve the project memory to answer.
func (m ProjectMemory) Index() string {
	return fmt.Sprintf("- [%s 发布约定](%s) — 项目的发布窗口、异常联系人与回滚约定。\n", m.Project, ProjectMemoryPath)
}

func (m ProjectMemory) Prompt() string {
	return fmt.Sprintf("请根据你记得的项目约定，为「%s」拟一份简短的上线安排，涵盖开始时间、异常联系和回滚处理。只需给出计划，不要执行发布；如果缺少信息，请明确说明。", m.Project)
}

func (m ProjectMemory) Verify(r *Run, result TurnResult) error {
	r.Step("检查上线安排是否用到了预先写入的记忆")
	if err := result.Verify(nil, false); err != nil {
		return err
	}
	var missing []string
	for _, fact := range []struct{ label, value string }{
		{"发布开始时间", m.ReleaseTime},
		{"异常联系人", m.Contact},
		{"回滚版本", m.RollbackVersion},
	} {
		if strings.Contains(result.Text, fact.value) {
			r.Log("checked", fact.label+"："+fact.value)
		} else {
			r.Log("info", "回复未体现"+fact.label+"（记忆中的值："+fact.value+"）")
			missing = append(missing, fact.label)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("助手的最终回复未体现以下记忆：%s；last_event_id=%s", strings.Join(missing, "、"), result.LastID)
	}
	r.Log("checked", "全新会话的回答体现了 3 项记忆；这些值只通过 Memory Store 提供")
	return nil
}
