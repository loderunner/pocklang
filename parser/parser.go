package parser

// Syntactical grammar:
// Expr       -> Or ;
// Or         -> And ("||" And)* ;
// And        -> Comparison ("&&" Comparison)* ;
// Comparison -> Term (("<" | ">" | ">=" | "<=" | "==" | "!=") Term) ;
// Term       -> Factor (("+" | "-") Factor)* ;
// Factor     -> Unary (("*" | "/") Unary)* ;
// Unary      -> ("!" | "-") Primary ;
// Primary    -> "true" | "false" | "null" | INTEGER | DECIMAL | STRING | "(" Expression ")" | IDENTIFIER ("." IDENTIFIER)* ;
