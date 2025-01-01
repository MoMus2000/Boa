from enum import Enum, auto

class TokenType(Enum):
    LEFT_PAREN           = auto()
    RIGHT_PAREN          = auto()
    LEFT_BRACE           = auto() 
    RIGHT_BRACE          = auto() 

    LEFT_ANGLE_BRACKET   = auto()
    RIGHT_ANGLE_BRACKET  = auto()

    PIPE                 = auto()

    COMMA                = auto()
    DOT                  = auto()
    MINUS                = auto()
    PLUS                 = auto()
    SEMICOLON            = auto()
    SLASH                = auto()
    STAR                 = auto()

    BANG                 = auto()
    BANG_EQUAL           = auto()
    EQUAL                = auto()
    EQUAL_EQUAL          = auto()
    GREATER              = auto()
    GREATER_EQUAL        = auto()
    LESS                 = auto()
    LESS_EQUAL           = auto()

    IDENTIFIER           = auto()
    STRING               = auto()
    NUMBER               = auto()

    AND                  = auto()
    CLASS                = auto()
    ELSE                 = auto()
    FALSE                = auto()
    FUN                  = auto()
    FOR                  = auto()
    IF                   = auto()
    NIL                  = auto()
    OR                   = auto()
    PRINT                = auto()
    RETURN               = auto()
    SUPER                = auto()
    THIS                 = auto()
    TRUE                 = auto()
    VAR                  = auto()
    WHILE                = auto()
    IMPORT               = auto()

    EOF                  = auto()