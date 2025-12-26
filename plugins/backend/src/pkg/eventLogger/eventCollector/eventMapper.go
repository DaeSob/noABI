package eventCollector

import (
	"cia/api/preference"
	"cia/common/blockchain/types"
	abiMapper "cia/common/blockchain/utils/events/abiMapper"

	"github.com/mitchellh/mapstructure"
)

func compileEvents(_functions []string) abiMapper.TSigMapper {

	mapper := abiMapper.TSigMapper{}
	// release memory
	defer func() {
		mapper = nil
	}()

	for _, function := range _functions {
		ev := types.CompileEvent(function)
		if ev != nil {
			mapper[ev.EncodeSig().HexString()] = *ev
		}
	}
	return mapper
}

// V2.0.0 By XeN
func getEventMapper(_eventLoggerInfo *preference.TEvetLoggerInfo) abiMapper.TABIMappers {

	collections := _eventLoggerInfo.Collections
	mapper := abiMapper.TABIMappers{}
	// release memory
	defer func() {
		mapper = nil
	}()

	eventIngerfaces := preference.GetEventLoggerEventInterfaces()
	for _, v := range collections {
		var collection preference.TCollection
		mapstructure.Decode(v, &collection)
		if collection.Enable {
			infs := (*eventIngerfaces)[collection.Interface]
			if len(infs) > 0 {
				mapper.SetSigMapper(
					collection.Address,
					compileEvents(infs), //getEventMap(name),
				)
			}
		}
	}
	return mapper

}
