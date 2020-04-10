package main

func dispatch(f flags) ierrori {

	if len(f.new) > 0 {
		ie := newUnit(f.new)
		if ie != nil {
			return ie
		}
		if len(f.questionAs) > 0 {
			ie := useTemplate(f.new, f.questionAs, true, f.questionInline)
			if ie != nil {
				return ie
			}
		}
		if len(f.answerAs) > 0 {
			ie := useTemplate(f.new, f.answerAs, false, f.answerInline)
			if ie != nil {
				return ie
			}
		}
	} else if f.question {
		ie := openQOrA(true)
		if ie != nil {
			return ie
		}
	} else if f.answer {
		ie := openQOrA(false)
		if ie != nil {
			return ie
		}
	} else if f.yes {
		ie := yesOrNo(true)
		if ie != nil {
			return ie
		}
	} else if f.no {
		ie := yesOrNo(false)
		if ie != nil {
			return ie
		}
	} else if f.addTemplate && len(f.rest) > 1 {
		ie := addTemplate(f.rest[0], f.rest[1])
		if ie != nil {
			return ie
		}
	} else if len(f.deleteTemplate) > 0 {
		ie := deleteTemplate(f.deleteTemplate)
		if ie != nil {
			return ie
		}
	} else if f.listTemplates {
		ie := listTemplates()
		if ie != nil {
			return ie
		}
	} else if f.renameTemplate && len(f.rest) > 1 {
		ie := renameTemplate(f.rest[0], f.rest[1])
		if ie != nil {
			return ie
		}
	} else {
		return ierror{m: "Unknown flag combination"}
	}

	return nil
}
