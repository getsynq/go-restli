package types

import (
	"github.com/PapaCharlie/go-restli/codegen/utils"
	. "github.com/dave/jennifer/jen"
)

const EnumShouldUsePointer = utils.No

type Enum struct {
	NamedType
	Symbols     []string
	SymbolToDoc map[string]string
}

func (e *Enum) InnerTypes() utils.IdentifierSet {
	return nil
}

func (e *Enum) ShouldReference() utils.ShouldUsePointer {
	return EnumShouldUsePointer
}

func (e *Enum) GenerateCode() (def *Statement) {
	def = Empty()
	utils.AddWordWrappedComment(def, e.Doc).Line()
	def.Type().Id(e.Name).Add(utils.Enum).Line()

	unknownEnum := Code(Id("_" + e.SymbolIdentifier("unknown")))

	def.Const().DefsFunc(func(def *Group) {
		def.Add(unknownEnum).Op("=").Id(e.Name).Call(Iota())
		for _, symbol := range e.Symbols {
			def.Add(utils.AddWordWrappedComment(Empty(), e.SymbolToDoc[symbol]))
			def.Id(e.SymbolIdentifier(symbol))
		}
	}).Line()

	values := "_" + e.Name + "_values"
	def.Var().Id(values).Op("=").Map(String()).Id(e.Name).Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Lit(s)] = Id(e.SymbolIdentifier(s))
		}
	})).Line()

	strings := "_" + e.Name + "_strings"
	def.Var().Id(strings).Op("=").Map(Id(e.Name)).String().Values(DictFunc(func(dict Dict) {
		for _, s := range e.Symbols {
			dict[Id(e.SymbolIdentifier(s))] = Lit(s)
		}
	})).Line().Line()

	val, ok := Code(Id("val")), Code(Id("ok"))
	receiver := e.Receiver()
	getter := "Get" + e.Name + "FromString"
	getEnumString := func(def *Group) {
		def.List(val, ok).Op(":=").Id(strings).Index(Id(receiver))
	}

	cast := func(c Code) *Statement {
		return Add(utils.Enum).Call(c)
	}

	utils.AddFuncOnReceiver(def, receiver, e.Name, utils.IsUnknown, EnumShouldUsePointer).Params().Bool().BlockFunc(func(def *Group) {
		def.Return(cast(Id(receiver)).Dot(utils.IsUnknown).Call())
	}).Line().Line()

	AddEquals(def, receiver, e.Name, EnumShouldUsePointer, func(other Code, def *Group) {
		def.Return(cast(Id(receiver)).Dot(utils.Equals).Call(cast(other)))
	})

	AddCustomComputeHash(def, receiver, e.Name, EnumShouldUsePointer, func(def *Group) {
		def.Return(cast(Id(receiver)).Dot(utils.ComputeHash).Call())
	})

	def.Func().Id("All" + e.Name + "Values").Params().Index().Id(e.Name).BlockFunc(func(def *Group) {
		def.Return(Index().Id(e.Name).ValuesFunc(func(def *Group) {
			for _, s := range e.Symbols {
				def.Id(e.SymbolIdentifier(s))
			}
		}))
	}).Line().Line()

	def.Func().Id(getter).Params(Add(val).String()).Params(Id(receiver).Id(e.Name), Err().Error()).
		BlockFunc(func(def *Group) {
			def.List(Id(receiver), ok).Op(":=").Id(values).Index(val)
			def.If(Op("!").Add(ok)).BlockFunc(func(def *Group) {
				def.Err().Op("=").Op("&").Add(utils.UnknownEnumValue).Values(Dict{
					Id("Enum"):  Lit(e.Identifier.String()),
					Id("Value"): val,
				})
			})
			def.Return(Id(receiver), Err())
		}).Line().Line()

	utils.AddStringer(def, receiver, e.Name, EnumShouldUsePointer, func(def *Group) {
		getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Lit("$UNKNOWN$")),
		)
		def.Return(val)
	}).Line().Line()

	utils.AddPointer(def, receiver, e.Name)

	AddMarshalRestLi(def, receiver, e.Name, EnumShouldUsePointer, func(def *Group) {
		getEnumString(def)
		def.If(Op("!").Add(ok)).Block(
			Return(Op("&").Add(utils.IllegalEnumConstant).Values(Dict{
				Id("Enum"):     Lit(e.Identifier.String()),
				Id("Constant"): Int().Call(Id(receiver)),
			})),
		)
		def.List(Writer).Dot("WriteString").Call(val)
		def.Return(Nil())
	})

	AddUnmarshalerFunc(def, receiver, e.Identifier, EnumShouldUsePointer)

	AddUnmarshalRestli(def, receiver, e.Name, func(def *Group) {
		value := Id("value")
		def.Var().Add(value).String()
		def.Add(Reader.Read(RestliType{Primitive: &StringPrimitive}, Reader, value))
		def.Add(utils.IfErrReturn(Err()))
		def.Line()

		def.Op("*").Id(receiver).Op("=").Id(values).Index(Id("value"))
		def.Return(Nil())
	})

	return def
}

func (e *Enum) UnmarshalerFunc() *Statement {
	return Qual(e.PackagePath(), e.unmarshalerFuncName())
}

func (e *Enum) unmarshalerFuncName() string {
	return "Unmarshal" + e.Name
}

func (e *Enum) SymbolIdentifier(symbol string) string {
	return utils.ExportedIdentifier(e.Name + "_" + symbol)
}

func (e *Enum) zeroValueLit() *Statement {
	return e.Qual().Call(Lit(0))
}

func (e *Enum) isValidSymbol(v string) bool {
	for _, symbol := range e.Symbols {
		if symbol == v {
			return true
		}
	}
	return false
}