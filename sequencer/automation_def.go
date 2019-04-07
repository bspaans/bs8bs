package sequencer

import "fmt"

type AutomationDef struct {
	BackAndForth *[]int           `yaml:"back_and_forth"`
	Cycle        *[]int           `yaml:"cycle"`
	Range        *RangeDef        `yaml:"range"`
	Sweep        *RangeDef        `yaml:"sweep"`
	FadeIn       *RangeDef        `yaml:"fade_in"`
	Register     *int             `yaml:"register"`
	Transpose    *IntTransposeDef `yaml:"transpose"`
}

func (a *AutomationDef) GetAutomation() (IntAutomation, error) {
	if a.BackAndForth != nil {
		return IntBackAndForthAutomation(*a.BackAndForth), nil
	} else if a.Cycle != nil {
		return IntCycleAutomation(*a.Cycle), nil
	} else if a.Range != nil {
		return IntRangeAutomation(a.Range.From, a.Range.To, a.Range.Step), nil
	} else if a.Sweep != nil {
		return IntSweepAutomation(a.Sweep.From, a.Sweep.To, a.Sweep.ChangeEvery), nil
	} else if a.FadeIn != nil {
		return IntFadeInAutomation(a.FadeIn.From, a.FadeIn.To, a.FadeIn.ChangeEvery), nil
	} else if a.Register != nil {
		return IntRegisterAutomation(*a.Register), nil
	} else if a.Transpose != nil {
		automation, err := a.Transpose.Automation.GetAutomation()
		if err != nil {
			return nil, err
		}
		return IntTransposeAutomation(a.Transpose.Transpose, automation), nil
	}
	return nil, fmt.Errorf("Missing automation")
}

type FloatAutomationDef struct {
	BackAndForth *[]float64         `yaml:"back_and_forth"`
	Register     *int               `yaml:"register"`
	Transpose    *FloatTransposeDef `yaml:"transpose"`
}

func (a *FloatAutomationDef) GetAutomation() (FloatAutomation, error) {
	if a.BackAndForth != nil {
		return FloatBackAndForthAutomation(*a.BackAndForth), nil
	} else if a.Register != nil {
		return FloatRegisterAutomation(*a.Register), nil
	} else if a.Transpose != nil {
		automation, err := a.Transpose.Automation.GetAutomation()
		if err != nil {
			return nil, err
		}
		return FloatTransposeAutomation(a.Transpose.Transpose, automation), nil
	}
	return nil, fmt.Errorf("Missing automation")
}

type RangeDef struct {
	From        int
	To          int
	Step        int
	ChangeEvery int `yaml:"change_every"`
}

type CycleChordsDef struct {
	Count  int     `yaml:"count"`
	Chords [][]int `yaml:"chords"`
}

func (c *CycleChordsDef) GetAutomation(seq *Sequencer) (IntArrayAutomation, error) {
	return ChordCycleArrayAutomation(c.Count, c.Chords), nil
}

type IntArrayAutomationDef struct {
	CycleChords *CycleChordsDef       `yaml:"cycle_chords"`
	Register    *int                  `yaml:"register"`
	Transpose   *IntArrayTransposeDef `yaml:"transpose"`
	Index       *IntArrayIndexDef     `yaml:"index"`
}

func (a *IntArrayAutomationDef) GetAutomation(seq *Sequencer) (IntArrayAutomation, error) {
	if a.Register != nil {
		return IntArrayRegisterAutomation(*a.Register), nil
	} else if a.CycleChords != nil {
		return a.CycleChords.GetAutomation(seq)
	} else if a.Transpose != nil {
		automation, err := a.Transpose.Automation.GetAutomation(seq)
		if err != nil {
			return nil, err
		}
		return IntArrayTransposeAutomation(a.Transpose.Transpose, automation), nil
	} else if a.Index != nil {
		automation, err := a.Index.Automation.GetAutomation(seq)
		if err != nil {
			return nil, err
		}
		return IntArrayIndexAutomation(a.Index.Index, automation), nil
	}
	return nil, fmt.Errorf("Missing array automation")
}

type IntArrayIndexDef struct {
	Index      int                   `yaml:"value"`
	Automation IntArrayAutomationDef `yaml:",inline"`
}

type IntTransposeDef struct {
	Transpose  int           `yaml:"value"`
	Automation AutomationDef `yaml:",inline"`
}
type FloatTransposeDef struct {
	Transpose  float64            `yaml:"value"`
	Automation FloatAutomationDef `yaml:",inline"`
}
type IntArrayTransposeDef struct {
	Transpose  int                   `yaml:"value"`
	Automation IntArrayAutomationDef `yaml:",inline"`
}