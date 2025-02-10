package main

import (
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

type Precidence uint

const COUNT_MAX = math.MaxInt

const DEBUG_PRINT_CODE = 0

const (
	PREC_NONE       Precidence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

func (c *Compiler) parsePrecidence(precidence Precidence) {
	c.advance()
	prefixRule := c.getRule(c.parser.previous.tokenType).prefix
	if prefixRule == nil {
		c.error("Expected expression")
		return
	}
	canAssign := precidence <= PREC_ASSIGNMENT
	prefixRule(canAssign)
	for precidence <= c.getRule(c.parser.current.tokenType).precidence {
		c.advance()
		infixRule := c.getRule(c.parser.previous.tokenType).infix
		infixRule(canAssign)
	}
}

type ParseRule struct {
	prefix     func(canAssign bool)
	infix      func(canAssign bool)
	precidence Precidence
}

type Parser struct {
	current   Token
	previous  Token
	hadError  bool
	panicMode bool
}

type FunctionType uint8

const (
	TYPE_SCRIPT FunctionType = iota
	TYPE_FUNC
)

type Compiler struct {
	enclosing    *Compiler
	scanner      *Scanner
	parser       *Parser
	parseRules   map[TokenType]ParseRule
	localCount   int
	scopeDepth   int
	locals       []Local
	functionType FunctionType
	function     *ObjectFunc
}

type Local struct {
	name  Token
	depth int
}

func NewCompiler() Compiler {
	return Compiler{}
}

var compilingChunk *Chunk

func (c *Compiler) buildParseRules() map[TokenType]ParseRule {
	return map[TokenType]ParseRule{
		LEFT_PAREN:    {c.grouping, nil, PREC_NONE},
		RIGHT_PAREN:   {nil, nil, PREC_NONE},
		LEFT_BRACE:    {nil, nil, PREC_NONE},
		RIGHT_BRACE:   {nil, nil, PREC_NONE},
		COMMA:         {nil, nil, PREC_NONE},
		DOT:           {nil, nil, PREC_NONE},
		MINUS:         {c.unary, c.binary, PREC_TERM},
		PLUS:          {nil, c.binary, PREC_TERM},
		SEMICOLON:     {nil, nil, PREC_NONE},
		SLASH:         {nil, c.binary, PREC_FACTOR},
		STAR:          {nil, c.binary, PREC_FACTOR},
		BANG:          {c.unary, nil, PREC_NONE},
		BANG_EQUAL:    {nil, nil, PREC_NONE},
		EQUAL:         {nil, nil, PREC_NONE},
		EQUAL_EQUAL:   {nil, c.binary, PREC_COMPARISON},
		GREATER:       {nil, c.binary, PREC_COMPARISON},
		GREATER_EQUAL: {nil, c.binary, PREC_COMPARISON},
		LESS:          {nil, c.binary, PREC_COMPARISON},
		LESS_EQUAL:    {nil, c.binary, PREC_COMPARISON},
		IDENTIFIER:    {c.variable, nil, PREC_NONE},
		STRING:        {c.str, nil, PREC_NONE},
		NUMBER:        {c.number, nil, PREC_NONE},
		AND:           {nil, c.and_, PREC_AND},
		CLASS:         {nil, nil, PREC_NONE},
		ELSE:          {nil, nil, PREC_NONE},
		FALSE:         {c.literal, nil, PREC_NONE},
		FOR:           {nil, nil, PREC_NONE},
		FUN:           {nil, nil, PREC_NONE},
		IF:            {nil, nil, PREC_NONE},
		NIL:           {c.literal, nil, PREC_NONE},
		OR:            {nil, c.or_, PREC_OR},
		PRINT:         {nil, nil, PREC_NONE},
		RETURN:        {nil, nil, PREC_NONE},
		SUPER:         {nil, nil, PREC_NONE},
		THIS:          {nil, nil, PREC_NONE},
		TRUE:          {c.literal, nil, PREC_NONE},
		VAR:           {nil, nil, PREC_NONE},
		WHILE:         {nil, nil, PREC_NONE},
		ERROR:         {nil, nil, PREC_NONE},
		EOF:           {nil, nil, PREC_NONE},
	}
}

func (c *Compiler) initNestedCompiler(compiler *Compiler, funcType FunctionType) {
	compiler.enclosing = c
	compiler.function = NewFunction()
	compiler.functionType = funcType
	compiler.parser = c.parser
	c = compiler
	if funcType != TYPE_SCRIPT {
		s := string(c.parser.previous.runes)
		objS := ObjectString{
			obj:    Object{objType: OBJ_STRING},
			chars:  s,
			length: len(s),
		}
		c.function.name = &objS
	}
	c.localCount++
	local := Local{}
	local.depth = 0
	local.name.runes = []rune("")
	local.name.length = 0
	c.locals = append(c.locals, local)
}

func (c *Compiler) compile(source []byte) *ObjectFunc {
	// ----------------------------
	parser := Parser{}
	scanner := NewScanner(source)
	c.parseRules = c.buildParseRules()
	c.scanner = &scanner
	c.parser = &parser
	c.localCount = 0
	c.scopeDepth = 0
	c.locals = make([]Local, 0)
	c.parser.hadError = false
	c.parser.panicMode = false
	c.function = NewFunction()
	c.functionType = TYPE_SCRIPT
	// ----------------------------
	c.localCount++
	local := Local{}
	local.depth = 0
	local.name.runes = []rune("")
	local.name.length = 0
	c.locals = append(c.locals, local)
	// ----------------------------
	c.advance()
	for !c.match(EOF) {
		c.declaration()
	}
	error := !c.parser.hadError
	if error == false {
		return nil
	}
	return c.endCompiler()
}

func (c *Compiler) declaration() {
	if c.match(FUN) {
		c.funDeclaration()
	} else if c.match(VAR) {
		c.varDeclaration()
	} else {
		c.statement()
	}
}

func (c *Compiler) funDeclaration() {
	global := c.parseVariable("Expect function name")
	c.markInitialized()
	c.makeFunction(TYPE_FUNC)
	c.defineVariable(global)
}

func (c *Compiler) makeFunction(funcType FunctionType) {
	compiler := NewCompiler()
	c.initNestedCompiler(&compiler, funcType)

	c.beginScope()
	c.consume(LEFT_PAREN, "Expected (")
	if !c.check(RIGHT_PAREN) {
		for {
			c.function.arity++
			if c.function.arity > 255 {
				c.error("Cant have more than 255 args.")
			}
			constant := c.parseVariable("Expect parameter name")
			c.defineVariable(constant)
			if !c.match(COMMA) {
				break
			}
		}
	}
	c.consume(RIGHT_PAREN, "Expected )")
	c.consume(LEFT_BRACE, "Expected {")
	c.block()

	// function := c.endCompiler()
	// objFunc := (*Object)(unsafe.Pointer(&function))
	// c.emitBytes(OpConstant, c.makeConstant(ObjVal(objFunc)))
}

func (c *Compiler) varDeclaration() {
	index := c.parseVariable("Expected ident")
	if c.match(EQUAL) {
		c.expression()
	} else {
		c.emitByteCode(OpNil)
	}
	c.consume(SEMICOLON, "Expected ;")
	c.defineVariable(index)
}

func (c *Compiler) markInitialized() {
	if c.scopeDepth == 0 {
		return
	}
	c.locals[c.localCount-1].depth = c.scopeDepth
}

func (c *Compiler) defineVariable(i Opcode) {
	if c.scopeDepth > 0 {
		c.markInitialized()
		return
	}
	c.emitBytes(OpDefineGlobal, i)
}

func (c *Compiler) parseVariable(msg string) Opcode {
	c.consume(IDENTIFIER, msg)
	c.declareVariable()
	if c.scopeDepth > 0 {
		return 0
	}
	ident := string(c.parser.previous.runes)
	str := ObjectString{
		obj:    Object{objType: OBJ_STRING},
		length: len(ident),
		chars:  ident,
	}
	obj := (*Object)(unsafe.Pointer(&str))
	index := c.makeConstant(
		ObjVal(obj),
	)
	return index
}

func (c *Compiler) declareVariable() {
	if c.scopeDepth == 0 {
		return
	}
	name := c.parser.previous
	i := c.localCount - 1
	for i >= 0 {
		local := c.locals[i]
		if local.depth != -1 && local.depth < c.scopeDepth {
			break
		}
		if c.identifierEquals(local.name, name) {
			panic(fmt.Sprint("Var with ", name, " Already exists"))
		}
		i--
	}
	c.addLocal(name)
}

func (c *Compiler) identifierEquals(t Token, u Token) bool {
	return string(t.runes) == string(u.runes) && t.tokenType == u.tokenType
}

func (c *Compiler) addLocal(name Token) {
	if len(c.locals) > COUNT_MAX {
		panic("Too many vars declared")
	}
	l := Local{
		name:  name,
		depth: -1,
	}
	c.locals = append(c.locals, l)
	c.localCount++
}

func (c *Compiler) statement() {
	if c.match(PRINT) {
		fmt.Println("Print Statement")
		c.printStatement()
		fmt.Println("Done with Print Statement")
	} else if c.match(LEFT_BRACE) {
		c.beginScope()
		c.block()
		c.endScope()
	} else if c.match(IF) {
		c.ifStatement()
	} else if c.match(WHILE) {
		c.whileStatement()
	} else {
		fmt.Println("Inside Expression Statement")
		c.expressionStatement()
	}
}

func (c *Compiler) whileStatement() {
	loopStart := len(c.currentChunk().code)
	c.consume(LEFT_PAREN, "Expected ) Paren")
	c.expression()
	c.consume(RIGHT_PAREN, "Expected ( Paren")

	exitJump := c.emitJump(OpJumpIfFalse)
	c.emitByteCode(OpPop)
	c.statement()
	c.emitLoop(loopStart)
	c.patchJump(exitJump)
	c.emitByteCode(OpPop)
}

func (c *Compiler) emitLoop(loopStart int) {
	c.emitByteCode(OpLoop)

	offset := len(c.currentChunk().code) - loopStart + 2

	if offset > math.MaxUint16 {
		c.error("Offset is too large.")
	}

	c.emitByteCode(Opcode(uint16(offset>>8) & 0xff))
	c.emitByteCode(Opcode(offset & 0xff))
}

func (c *Compiler) ifStatement() {
	c.consume(LEFT_PAREN, "Expected (")
	c.expression()
	c.consume(RIGHT_PAREN, "Expected )")

	thenJump := c.emitJump(OpJumpIfFalse)
	c.statement()
	c.emitByteCode(OpPop)
	elseJump := c.emitJump(OpJump)
	c.patchJump(thenJump)

	c.emitByteCode(OpPop)
	if c.match(ELSE) {
		c.statement()
	}
	c.patchJump(elseJump)
}

func (c *Compiler) emitJump(jCode Opcode) int {
	c.emitByteCode(jCode)
	c.emitByteCode(0xff)
	c.emitByteCode(0xff)
	return len(c.currentChunk().code) - 2
}

func (c *Compiler) patchJump(offset int) {
	jump := len(c.currentChunk().code) - offset - 2
	if jump > math.MaxUint16 {
		c.error("Too much code to jump over.")
	}
	c.currentChunk().code[offset] = Opcode((jump >> 8) & 0xff)
	c.currentChunk().code[offset+1] = Opcode(jump & 0xff)
}

func (c *Compiler) and_(canAssign bool) {
	endJump := c.emitJump(OpJumpIfFalse)
	c.emitByteCode(OpPop)
	c.parsePrecidence(PREC_AND)
	c.patchJump(endJump)
}

func (c *Compiler) or_(canAssign bool) {
	elseJump := c.emitJump(OpJumpIfFalse)
	endJump := c.emitJump(OpJump)
	c.patchJump(elseJump)
	c.emitByteCode(OpPop)
	c.parsePrecidence(PREC_OR)
	c.patchJump(endJump)
}

func (c *Compiler) beginScope() {
	c.scopeDepth++
}

func (c *Compiler) endScope() {
	c.scopeDepth--
	for c.localCount > 0 && c.locals[c.localCount-1].depth > c.scopeDepth {
		c.emitByteCode(OpPop)
		c.localCount--
	}
}

func (c *Compiler) block() {
	for !c.check(RIGHT_BRACE) && !c.check(EOF) {
		c.declaration()
	}
	c.consume(RIGHT_BRACE, "Expected }")
}

func (c *Compiler) printStatement() {
	c.expression()
	fmt.Println("PArsed Expression")
	c.consume(SEMICOLON, "Expected ;")
	c.emitByteCode(OpPrint)
}

func (c *Compiler) expressionStatement() {
	c.expression()
	c.consume(SEMICOLON, "Expected ;")
	c.emitByteCode(OpPop)
}

func (c *Compiler) match(token TokenType) bool {
	if !c.check(token) {
		return false
	}
	c.advance()
	return true
}

func (c *Compiler) check(token TokenType) bool {
	return c.parser.current.tokenType == token
}

func (c *Compiler) endCompiler() *ObjectFunc {
	if !c.parser.hadError && DEBUG_PRINT_CODE == 1 {
		DisassembleChunk(c.currentChunk(), "code")
	}
	c.emitReturn()
	function := c.function
	c = c.enclosing
	return function
}

func (c *Compiler) variable(canAssign bool) {
	var getOp, setOp Opcode
	ident := string(c.parser.previous.runes)
	arg := c.resolveLocal(c.parser.previous)
	if arg != OpMinus1 {
		getOp = OpGetLocal
		setOp = OpSetLocal
	} else {
		objectString := ObjectString{
			obj:    Object{objType: OBJ_STRING},
			length: len(ident),
			chars:  ident,
		}
		arg = c.makeConstant(ObjVal((*Object)(unsafe.Pointer(&objectString))))
		getOp = OpGetGlobal
		setOp = OpSetGlobal
	}
	if canAssign && c.match(EQUAL) {
		c.expression()
		c.emitBytes(setOp, arg)
	} else {
		c.emitBytes(getOp, arg)
	}
}

func (c *Compiler) resolveLocal(name Token) Opcode {
	i := len(c.locals) - 1
	for i >= 0 {
		local := c.locals[i]
		if c.identifierEquals(name, local.name) {
			return Opcode(i)
		}
		i--
	}
	return OpMinus1
}

func (c *Compiler) number(canAssign bool) {
	num, err := strconv.ParseFloat(string(c.parser.previous.runes), 32)
	if err != nil {
		err := err.Error()
		c.errorAtCurrent(err)
	}
	c.emitBytes(OpConstant, c.makeConstant(NumberVal(float32(num))))
}

func (c *Compiler) grouping(canAssign bool) {
	c.expression()
	c.consume(RIGHT_PAREN, fmt.Sprintf("Expected ) after expression, but got : %s", c.parser.current.tokenType))
}

func (c *Compiler) unary(canAssign bool) {
	operator := c.parser.previous.tokenType
	//c.expression() // Evaluate the operand first and then apply whatever operator we have to
	c.parsePrecidence(PREC_UNARY)
	switch operator {
	case MINUS:
		c.emitByteCode(OpNegate)
		break
	case BANG:
		c.emitByteCode(OpNot)
		break
	default:
		return
	}
}

func (c *Compiler) binary(canAssign bool) {
	operator := c.parser.previous.tokenType
	rule := c.getRule(operator)
	c.parsePrecidence(Precidence(rule.precidence + 1))

	switch operator {
	case PLUS:
		c.emitByteCode(OpAdd)
		break
	case MINUS:
		c.emitByteCode(OpSub)
		break
	case STAR:
		c.emitByteCode(OpMul)
		break
	case SLASH:
		c.emitByteCode(OpDiv)
		break
	case GREATER:
		c.emitByteCode(OpGreater)
		break
	case LESS:
		c.emitByteCode(OpLess)
		break
	case EQUAL_EQUAL:
		c.emitByteCode(OpEqual)
		break
	case LESS_EQUAL:
		c.emitBytes(OpGreater, OpNot)
		break
	case GREATER_EQUAL:
		c.emitBytes(OpLess, OpNot)
		break
	default:
		return
	}
}

func (c *Compiler) literal(canAssign bool) {
	op := c.parser.previous.tokenType
	switch op {
	case TRUE:
		c.emitByteCode(OpTrue)
	case FALSE:
		c.emitByteCode(OpFalse)
	case NIL:
		c.emitByteCode(OpNil)
	default:
		return
	}
}

func (c *Compiler) str(canAssign bool) {
	str := string(c.parser.previous.runes)
	objectString := ObjectString{
		obj:    Object{objType: OBJ_STRING},
		length: len(str),
		chars:  str,
	}
	object := (*Object)(unsafe.Pointer(&objectString))
	c.emitBytes(OpConstant, c.makeConstant(ObjVal(object)))
}

func (c *Compiler) getRule(token TokenType) *ParseRule {
	pRule := c.parseRules[token]
	return &pRule
}

func (c *Compiler) makeConstant(constant Value) Opcode {
	index := c.currentChunk().AddConstant(constant)
	if index > int(^uint(0)>>1) {
		c.error("Too many constants in the chunk")
		return 0
	}
	return Opcode(index)
}

func (c *Compiler) emitReturn() {
	c.emitByteCode(OpReturn)
}

func (c *Compiler) emitBytes(a Opcode, b Opcode) {
	c.emitByteCode(a)
	c.emitByteCode(b)
}

func (c *Compiler) currentChunk() *Chunk {
	return &c.function.chunk
}

func (c *Compiler) expression() {
	c.parsePrecidence(PREC_ASSIGNMENT)
}

func (c *Compiler) advance() {
	c.parser.previous = c.parser.current
	for {
		token := c.scanner.scanToken()
		c.parser.current = token
		if c.parser.current.tokenType != ERROR {
			break
		}
		c.errorAtCurrent(string(c.parser.current.runes))
	}
}

func (c *Compiler) consume(tokenType TokenType, message string) {
	if c.parser.current.tokenType == tokenType {
		c.advance()
		return
	}
	c.errorAtCurrent(message)
}

func (c *Compiler) emitByteCode(code Opcode) {
	c.currentChunk().WriteChunk(code, c.parser.previous.line)
}

func (c *Compiler) errorAtCurrent(message string) {
	c.errorAt(&c.parser.current, message)
}

func (c *Compiler) error(message string) {
	c.errorAt(&c.parser.previous, message)
}

func (c *Compiler) errorAt(token *Token, message string) {
	if c.parser.panicMode {
		return
	}
	c.parser.panicMode = true
	fmt.Printf("[line %d] Error", token.line)
	if token.tokenType == EOF {
		fmt.Printf(" at end")
	} else if token.tokenType == ERROR {

	} else {
		fmt.Printf(" at '%.*s'", token.length, string(token.runes))
	}
	fmt.Printf(": %s\n", message)
	c.parser.hadError = true
}
