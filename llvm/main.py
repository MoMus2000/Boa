from llvmlite import ir
import llvmlite.binding as llvm

def compile():
    module  = ir.Module('main')
    return_type = ir.IntType(32)
    function_type = ir.FunctionType(return_type, [])
    func = ir.Function(module, function_type, name="main")
    block = func.append_basic_block("main_entry")
    builder = ir.IRBuilder(block)
    return_value = ir.Constant(ir.IntType(32), 420)
    builder.ret(return_value)

    # get the architecture to generate code for
    module.triple = llvm.get_default_triple()

    with open("./target/main.ll", "w") as f:
        f.write(str(module))

compile()

