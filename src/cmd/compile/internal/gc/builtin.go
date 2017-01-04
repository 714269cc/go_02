// AUTO-GENERATED by mkbuiltin.go; DO NOT EDIT

package gc

var runtimeDecls = [...]struct {
	name string
	tag  int
	typ  int
}{
	{"newobject", funcTag, 4},
	{"panicindex", funcTag, 5},
	{"panicslice", funcTag, 5},
	{"panicdivide", funcTag, 5},
	{"throwinit", funcTag, 5},
	{"panicwrap", funcTag, 5},
	{"gopanic", funcTag, 7},
	{"gorecover", funcTag, 10},
	{"goschedguarded", funcTag, 5},
	{"printbool", funcTag, 12},
	{"printfloat", funcTag, 14},
	{"printint", funcTag, 16},
	{"printhex", funcTag, 18},
	{"printuint", funcTag, 18},
	{"printcomplex", funcTag, 20},
	{"printstring", funcTag, 22},
	{"printpointer", funcTag, 23},
	{"printiface", funcTag, 23},
	{"printeface", funcTag, 23},
	{"printslice", funcTag, 23},
	{"printnl", funcTag, 5},
	{"printsp", funcTag, 5},
	{"printlock", funcTag, 5},
	{"printunlock", funcTag, 5},
	{"concatstring2", funcTag, 26},
	{"concatstring3", funcTag, 27},
	{"concatstring4", funcTag, 28},
	{"concatstring5", funcTag, 29},
	{"concatstrings", funcTag, 31},
	{"cmpstring", funcTag, 33},
	{"eqstring", funcTag, 34},
	{"intstring", funcTag, 37},
	{"slicebytetostring", funcTag, 39},
	{"slicebytetostringtmp", funcTag, 40},
	{"slicerunetostring", funcTag, 43},
	{"stringtoslicebyte", funcTag, 44},
	{"stringtoslicerune", funcTag, 47},
	{"decoderune", funcTag, 48},
	{"slicecopy", funcTag, 50},
	{"slicestringcopy", funcTag, 51},
	{"convI2I", funcTag, 52},
	{"convT2E", funcTag, 53},
	{"convT2I", funcTag, 53},
	{"assertE2I", funcTag, 52},
	{"assertE2I2", funcTag, 54},
	{"assertI2I", funcTag, 52},
	{"assertI2I2", funcTag, 54},
	{"panicdottypeE", funcTag, 55},
	{"panicdottypeI", funcTag, 55},
	{"panicnildottype", funcTag, 56},
	{"ifaceeq", funcTag, 57},
	{"efaceeq", funcTag, 57},
	{"makemap", funcTag, 59},
	{"mapaccess1", funcTag, 60},
	{"mapaccess1_fast32", funcTag, 61},
	{"mapaccess1_fast64", funcTag, 61},
	{"mapaccess1_faststr", funcTag, 61},
	{"mapaccess1_fat", funcTag, 62},
	{"mapaccess2", funcTag, 63},
	{"mapaccess2_fast32", funcTag, 64},
	{"mapaccess2_fast64", funcTag, 64},
	{"mapaccess2_faststr", funcTag, 64},
	{"mapaccess2_fat", funcTag, 65},
	{"mapassign", funcTag, 60},
	{"mapiterinit", funcTag, 66},
	{"mapdelete", funcTag, 66},
	{"mapiternext", funcTag, 67},
	{"makechan", funcTag, 69},
	{"chanrecv1", funcTag, 71},
	{"chanrecv2", funcTag, 72},
	{"chansend1", funcTag, 74},
	{"closechan", funcTag, 23},
	{"writeBarrier", varTag, 75},
	{"writebarrierptr", funcTag, 76},
	{"typedmemmove", funcTag, 77},
	{"typedmemclr", funcTag, 78},
	{"typedslicecopy", funcTag, 79},
	{"selectnbsend", funcTag, 80},
	{"selectnbrecv", funcTag, 81},
	{"selectnbrecv2", funcTag, 83},
	{"newselect", funcTag, 84},
	{"selectsend", funcTag, 80},
	{"selectrecv", funcTag, 72},
	{"selectrecv2", funcTag, 85},
	{"selectdefault", funcTag, 86},
	{"selectgo", funcTag, 56},
	{"block", funcTag, 5},
	{"makeslice", funcTag, 88},
	{"makeslice64", funcTag, 89},
	{"growslice", funcTag, 90},
	{"memmove", funcTag, 91},
	{"memclrNoHeapPointers", funcTag, 92},
	{"memclrHasPointers", funcTag, 92},
	{"memequal", funcTag, 93},
	{"memequal8", funcTag, 94},
	{"memequal16", funcTag, 94},
	{"memequal32", funcTag, 94},
	{"memequal64", funcTag, 94},
	{"memequal128", funcTag, 94},
	{"int64div", funcTag, 95},
	{"uint64div", funcTag, 96},
	{"int64mod", funcTag, 95},
	{"uint64mod", funcTag, 96},
	{"float64toint64", funcTag, 97},
	{"float64touint64", funcTag, 98},
	{"float64touint32", funcTag, 100},
	{"int64tofloat64", funcTag, 101},
	{"uint64tofloat64", funcTag, 102},
	{"uint32tofloat64", funcTag, 103},
	{"complex128div", funcTag, 104},
	{"racefuncenter", funcTag, 105},
	{"racefuncexit", funcTag, 5},
	{"raceread", funcTag, 105},
	{"racewrite", funcTag, 105},
	{"racereadrange", funcTag, 106},
	{"racewriterange", funcTag, 106},
	{"msanread", funcTag, 106},
	{"msanwrite", funcTag, 106},
}

func runtimeTypes() []*Type {
	var typs [107]*Type
	typs[0] = bytetype
	typs[1] = typPtr(typs[0])
	typs[2] = Types[TANY]
	typs[3] = typPtr(typs[2])
	typs[4] = functype(nil, []*Node{anonfield(typs[1])}, []*Node{anonfield(typs[3])})
	typs[5] = functype(nil, nil, nil)
	typs[6] = Types[TINTER]
	typs[7] = functype(nil, []*Node{anonfield(typs[6])}, nil)
	typs[8] = Types[TINT32]
	typs[9] = typPtr(typs[8])
	typs[10] = functype(nil, []*Node{anonfield(typs[9])}, []*Node{anonfield(typs[6])})
	typs[11] = Types[TBOOL]
	typs[12] = functype(nil, []*Node{anonfield(typs[11])}, nil)
	typs[13] = Types[TFLOAT64]
	typs[14] = functype(nil, []*Node{anonfield(typs[13])}, nil)
	typs[15] = Types[TINT64]
	typs[16] = functype(nil, []*Node{anonfield(typs[15])}, nil)
	typs[17] = Types[TUINT64]
	typs[18] = functype(nil, []*Node{anonfield(typs[17])}, nil)
	typs[19] = Types[TCOMPLEX128]
	typs[20] = functype(nil, []*Node{anonfield(typs[19])}, nil)
	typs[21] = Types[TSTRING]
	typs[22] = functype(nil, []*Node{anonfield(typs[21])}, nil)
	typs[23] = functype(nil, []*Node{anonfield(typs[2])}, nil)
	typs[24] = typArray(typs[0], 32)
	typs[25] = typPtr(typs[24])
	typs[26] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[27] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[28] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[29] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[30] = typSlice(typs[21])
	typs[31] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[30])}, []*Node{anonfield(typs[21])})
	typs[32] = Types[TINT]
	typs[33] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[32])})
	typs[34] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[11])})
	typs[35] = typArray(typs[0], 4)
	typs[36] = typPtr(typs[35])
	typs[37] = functype(nil, []*Node{anonfield(typs[36]), anonfield(typs[15])}, []*Node{anonfield(typs[21])})
	typs[38] = typSlice(typs[0])
	typs[39] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[38])}, []*Node{anonfield(typs[21])})
	typs[40] = functype(nil, []*Node{anonfield(typs[38])}, []*Node{anonfield(typs[21])})
	typs[41] = runetype
	typs[42] = typSlice(typs[41])
	typs[43] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[42])}, []*Node{anonfield(typs[21])})
	typs[44] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21])}, []*Node{anonfield(typs[38])})
	typs[45] = typArray(typs[41], 32)
	typs[46] = typPtr(typs[45])
	typs[47] = functype(nil, []*Node{anonfield(typs[46]), anonfield(typs[21])}, []*Node{anonfield(typs[42])})
	typs[48] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[32])}, []*Node{anonfield(typs[41]), anonfield(typs[32])})
	typs[49] = Types[TUINTPTR]
	typs[50] = functype(nil, []*Node{anonfield(typs[2]), anonfield(typs[2]), anonfield(typs[49])}, []*Node{anonfield(typs[32])})
	typs[51] = functype(nil, []*Node{anonfield(typs[2]), anonfield(typs[2])}, []*Node{anonfield(typs[32])})
	typs[52] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2])}, []*Node{anonfield(typs[2])})
	typs[53] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3])}, []*Node{anonfield(typs[2])})
	typs[54] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2])}, []*Node{anonfield(typs[2]), anonfield(typs[11])})
	typs[55] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[1]), anonfield(typs[1])}, nil)
	typs[56] = functype(nil, []*Node{anonfield(typs[1])}, nil)
	typs[57] = functype(nil, []*Node{anonfield(typs[2]), anonfield(typs[2])}, []*Node{anonfield(typs[11])})
	typs[58] = typMap(typs[2], typs[2])
	typs[59] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[3]), anonfield(typs[3])}, []*Node{anonfield(typs[58])})
	typs[60] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[58]), anonfield(typs[3])}, []*Node{anonfield(typs[3])})
	typs[61] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[58]), anonfield(typs[2])}, []*Node{anonfield(typs[3])})
	typs[62] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[58]), anonfield(typs[3]), anonfield(typs[1])}, []*Node{anonfield(typs[3])})
	typs[63] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[58]), anonfield(typs[3])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[64] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[58]), anonfield(typs[2])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[65] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[58]), anonfield(typs[3]), anonfield(typs[1])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[66] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[58]), anonfield(typs[3])}, nil)
	typs[67] = functype(nil, []*Node{anonfield(typs[3])}, nil)
	typs[68] = typChan(typs[2], Cboth)
	typs[69] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15])}, []*Node{anonfield(typs[68])})
	typs[70] = typChan(typs[2], Crecv)
	typs[71] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[70]), anonfield(typs[3])}, nil)
	typs[72] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[70]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[73] = typChan(typs[2], Csend)
	typs[74] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[73]), anonfield(typs[3])}, nil)
	typs[75] = tostruct([]*Node{namedfield("enabled", typs[11]), namedfield("needed", typs[11]), namedfield("cgo", typs[11])})
	typs[76] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[2])}, nil)
	typs[77] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3]), anonfield(typs[3])}, nil)
	typs[78] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3])}, nil)
	typs[79] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2]), anonfield(typs[2])}, []*Node{anonfield(typs[32])})
	typs[80] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[73]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[81] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3]), anonfield(typs[70])}, []*Node{anonfield(typs[11])})
	typs[82] = typPtr(typs[11])
	typs[83] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3]), anonfield(typs[82]), anonfield(typs[70])}, []*Node{anonfield(typs[11])})
	typs[84] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[8])}, nil)
	typs[85] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[70]), anonfield(typs[3]), anonfield(typs[82])}, []*Node{anonfield(typs[11])})
	typs[86] = functype(nil, []*Node{anonfield(typs[1])}, []*Node{anonfield(typs[11])})
	typs[87] = typSlice(typs[2])
	typs[88] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[32]), anonfield(typs[32])}, []*Node{anonfield(typs[87])})
	typs[89] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[15])}, []*Node{anonfield(typs[87])})
	typs[90] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[87]), anonfield(typs[32])}, []*Node{anonfield(typs[87])})
	typs[91] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3]), anonfield(typs[49])}, nil)
	typs[92] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[49])}, nil)
	typs[93] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3]), anonfield(typs[49])}, []*Node{anonfield(typs[11])})
	typs[94] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[95] = functype(nil, []*Node{anonfield(typs[15]), anonfield(typs[15])}, []*Node{anonfield(typs[15])})
	typs[96] = functype(nil, []*Node{anonfield(typs[17]), anonfield(typs[17])}, []*Node{anonfield(typs[17])})
	typs[97] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[15])})
	typs[98] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[17])})
	typs[99] = Types[TUINT32]
	typs[100] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[99])})
	typs[101] = functype(nil, []*Node{anonfield(typs[15])}, []*Node{anonfield(typs[13])})
	typs[102] = functype(nil, []*Node{anonfield(typs[17])}, []*Node{anonfield(typs[13])})
	typs[103] = functype(nil, []*Node{anonfield(typs[99])}, []*Node{anonfield(typs[13])})
	typs[104] = functype(nil, []*Node{anonfield(typs[19]), anonfield(typs[19])}, []*Node{anonfield(typs[19])})
	typs[105] = functype(nil, []*Node{anonfield(typs[49])}, nil)
	typs[106] = functype(nil, []*Node{anonfield(typs[49]), anonfield(typs[49])}, nil)
	return typs[:]
}
