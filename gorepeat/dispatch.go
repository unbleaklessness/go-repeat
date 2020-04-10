package main

func dispatch(f flags) ierrori {

	if len(f.new) > 0 {
		ie := newUnit(f.new, f.inverse)
		if ie != nil {
			return ie
		}
		ie = useTemplate(f.new, f.AAs, true, f.AInline)
		if ie != nil {
			return ie
		}
		ie = useTemplate(f.new, f.BAs, false, f.BInline)
		if ie != nil {
			return ie
		}
	} else if f.question {
		ie := openQuestionOrAnswer(true)
		if ie != nil {
			return ie
		}
	} else if f.answer {
		ie := openQuestionOrAnswer(false)
		if ie != nil {
			return ie
		}
	} else if f.yes {
		ie := answer(true)
		if ie != nil {
			return ie
		}
	} else if f.no {
		ie := answer(false)
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
	} else if f.inverse && len(f.rest) > 0 {
		ie := toggleInverse(f.rest[0])
		if ie != nil {
			return ie
		}
	} else if len(f.setDefaultInline) > 0 {
		ie := setDefaultInline(f.setDefaultInline)
		if ie != nil {
			return ie
		}
	} else {
		return ierror{m: "Unknown flag combination"}
	}

	return nil
}
