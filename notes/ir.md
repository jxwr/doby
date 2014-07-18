## IR

tmp = binary_add b c
a = tmp

bAddr = getAddr('b')
store(bAddr, 2)

off = getAddr('c')
store(off, 3)

aAddr = getAddr('a')

// tmp = b + c
bAddr = getAddr('b')
cAddr = getAddr('c')
tmpAddr = NewTmp()
store(tmpAddr, bAddr, cAddr)
return tmpAddr

// a = tmp
yAddr = build(Y)
aAddr = getAddr('a')
store(aAddr, yAddr)

// a.Times
addr = IntegerOperand(10)
tmp = NewTmp()
selector(addr, StringOperand('Times'), tmp)
return tmp

// fn(1,2,3)
addr = getMAddr('fn')

call(addr, )

n = AllocVar('b')
i = IntegerOperand(2)
instr.Copy(n, i)

n = AllocVar('c')
i = IntegerOperand(3)
instr.Copy(n, i)

buildIdent()
  return 

node {
     X Y
}

obj = build(X)
arg = build(Y)

instr.Call(obj, '+', arg)
instr.Copy()
