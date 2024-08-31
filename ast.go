package pock

type Expr interface{}

type BinaryExpr struct {
	Op    TokenType
	Left  Expr
	Right Expr
}

type UnaryExpr struct {
	Op   TokenType
	Expr Expr
}

type GroupExpr struct {
	Expr Expr
}

type GetExpr struct {
	Names []string
}

type LiteralExpr struct {
	Token Token
}
