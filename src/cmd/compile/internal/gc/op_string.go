// Code generated by "stringer -type=Op -trimprefix=O"; DO NOT EDIT.

package gc

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OXXX-0]
	_ = x[ONAME-1]
	_ = x[ONONAME-2]
	_ = x[OTYPE-3]
	_ = x[OPACK-4]
	_ = x[OLITERAL-5]
	_ = x[ONIL-6]
	_ = x[OADD-7]
	_ = x[OSUB-8]
	_ = x[OOR-9]
	_ = x[OXOR-10]
	_ = x[OADDSTR-11]
	_ = x[OADDR-12]
	_ = x[OANDAND-13]
	_ = x[OAPPEND-14]
	_ = x[OBYTES2STR-15]
	_ = x[OBYTES2STRTMP-16]
	_ = x[ORUNES2STR-17]
	_ = x[OSTR2BYTES-18]
	_ = x[OSTR2BYTESTMP-19]
	_ = x[OSTR2RUNES-20]
	_ = x[OAS-21]
	_ = x[OAS2-22]
	_ = x[OAS2DOTTYPE-23]
	_ = x[OAS2FUNC-24]
	_ = x[OAS2MAPR-25]
	_ = x[OAS2RECV-26]
	_ = x[OASOP-27]
	_ = x[OCALL-28]
	_ = x[OCALLFUNC-29]
	_ = x[OCALLMETH-30]
	_ = x[OCALLINTER-31]
	_ = x[OCALLPART-32]
	_ = x[OCAP-33]
	_ = x[OCLOSE-34]
	_ = x[OCLOSURE-35]
	_ = x[OCOMPLIT-36]
	_ = x[OMAPLIT-37]
	_ = x[OSTRUCTLIT-38]
	_ = x[OARRAYLIT-39]
	_ = x[OSLICELIT-40]
	_ = x[OPTRLIT-41]
	_ = x[OCONV-42]
	_ = x[OCONVIFACE-43]
	_ = x[OCONVNOP-44]
	_ = x[OCOPY-45]
	_ = x[ODCL-46]
	_ = x[ODCLFUNC-47]
	_ = x[ODCLFIELD-48]
	_ = x[ODCLCONST-49]
	_ = x[ODCLTYPE-50]
	_ = x[ODELETE-51]
	_ = x[ODOT-52]
	_ = x[ODOTPTR-53]
	_ = x[ODOTMETH-54]
	_ = x[ODOTINTER-55]
	_ = x[OXDOT-56]
	_ = x[ODOTTYPE-57]
	_ = x[ODOTTYPE2-58]
	_ = x[OEQ-59]
	_ = x[ONE-60]
	_ = x[OLT-61]
	_ = x[OLE-62]
	_ = x[OGE-63]
	_ = x[OGT-64]
	_ = x[ODEREF-65]
	_ = x[OINDEX-66]
	_ = x[OINDEXMAP-67]
	_ = x[OKEY-68]
	_ = x[OSTRUCTKEY-69]
	_ = x[OLEN-70]
	_ = x[OMAKE-71]
	_ = x[OMAKECHAN-72]
	_ = x[OMAKEMAP-73]
	_ = x[OMAKESLICE-74]
	_ = x[OMAKESLICECOPY-75]
	_ = x[OMUL-76]
	_ = x[ODIV-77]
	_ = x[OMOD-78]
	_ = x[OLSH-79]
	_ = x[ORSH-80]
	_ = x[OAND-81]
	_ = x[OANDNOT-82]
	_ = x[ONEW-83]
	_ = x[ONEWOBJ-84]
	_ = x[ONOT-85]
	_ = x[OBITNOT-86]
	_ = x[OPLUS-87]
	_ = x[ONEG-88]
	_ = x[OOROR-89]
	_ = x[OPANIC-90]
	_ = x[OPRINT-91]
	_ = x[OPRINTN-92]
	_ = x[OPAREN-93]
	_ = x[OSEND-94]
	_ = x[OSLICE-95]
	_ = x[OSLICEARR-96]
	_ = x[OSLICESTR-97]
	_ = x[OSLICE3-98]
	_ = x[OSLICE3ARR-99]
	_ = x[OSLICEHEADER-100]
	_ = x[ORECOVER-101]
	_ = x[ORECV-102]
	_ = x[ORUNESTR-103]
	_ = x[OSELRECV-104]
	_ = x[OSELRECV2-105]
	_ = x[OIOTA-106]
	_ = x[OREAL-107]
	_ = x[OIMAG-108]
	_ = x[OCOMPLEX-109]
	_ = x[OALIGNOF-110]
	_ = x[OOFFSETOF-111]
	_ = x[OSIZEOF-112]
	_ = x[OBLOCK-113]
	_ = x[OBREAK-114]
	_ = x[OCASE-115]
	_ = x[OCONTINUE-116]
	_ = x[ODEFER-117]
	_ = x[OEMPTY-118]
	_ = x[OFALL-119]
	_ = x[OFOR-120]
	_ = x[OFORUNTIL-121]
	_ = x[OGOTO-122]
	_ = x[OIF-123]
	_ = x[OLABEL-124]
	_ = x[OGO-125]
	_ = x[ORANGE-126]
	_ = x[ORETURN-127]
	_ = x[OSELECT-128]
	_ = x[OSWITCH-129]
	_ = x[OTYPESW-130]
	_ = x[OTCHAN-131]
	_ = x[OTMAP-132]
	_ = x[OTSTRUCT-133]
	_ = x[OTINTER-134]
	_ = x[OTFUNC-135]
	_ = x[OTARRAY-136]
	_ = x[ODDD-137]
	_ = x[OINLCALL-138]
	_ = x[OEFACE-139]
	_ = x[OITAB-140]
	_ = x[OIDATA-141]
	_ = x[OSPTR-142]
	_ = x[OCLOSUREVAR-143]
	_ = x[OCFUNC-144]
	_ = x[OCHECKNIL-145]
	_ = x[OVARDEF-146]
	_ = x[OVARKILL-147]
	_ = x[OVARLIVE-148]
	_ = x[ORESULT-149]
	_ = x[OINLMARK-150]
	_ = x[ORETJMP-151]
	_ = x[OGETG-152]
	_ = x[OEND-153]
}

const _Op_name = "XXXNAMENONAMETYPEPACKLITERALNILADDSUBORXORADDSTRADDRANDANDAPPENDBYTES2STRBYTES2STRTMPRUNES2STRSTR2BYTESSTR2BYTESTMPSTR2RUNESASAS2AS2DOTTYPEAS2FUNCAS2MAPRAS2RECVASOPCALLCALLFUNCCALLMETHCALLINTERCALLPARTCAPCLOSECLOSURECOMPLITMAPLITSTRUCTLITARRAYLITSLICELITPTRLITCONVCONVIFACECONVNOPCOPYDCLDCLFUNCDCLFIELDDCLCONSTDCLTYPEDELETEDOTDOTPTRDOTMETHDOTINTERXDOTDOTTYPEDOTTYPE2EQNELTLEGEGTDEREFINDEXINDEXMAPKEYSTRUCTKEYLENMAKEMAKECHANMAKEMAPMAKESLICEMAKESLICECOPYMULDIVMODLSHRSHANDANDNOTNEWNEWOBJNOTBITNOTPLUSNEGORORPANICPRINTPRINTNPARENSENDSLICESLICEARRSLICESTRSLICE3SLICE3ARRSLICEHEADERRECOVERRECVRUNESTRSELRECVSELRECV2IOTAREALIMAGCOMPLEXALIGNOFOFFSETOFSIZEOFBLOCKBREAKCASECONTINUEDEFEREMPTYFALLFORFORUNTILGOTOIFLABELGORANGERETURNSELECTSWITCHTYPESWTCHANTMAPTSTRUCTTINTERTFUNCTARRAYDDDINLCALLEFACEITABIDATASPTRCLOSUREVARCFUNCCHECKNILVARDEFVARKILLVARLIVERESULTINLMARKRETJMPGETGEND"

var _Op_index = [...]uint16{0, 3, 7, 13, 17, 21, 28, 31, 34, 37, 39, 42, 48, 52, 58, 64, 73, 85, 94, 103, 115, 124, 126, 129, 139, 146, 153, 160, 164, 168, 176, 184, 193, 201, 204, 209, 216, 223, 229, 238, 246, 254, 260, 264, 273, 280, 284, 287, 294, 302, 310, 317, 323, 326, 332, 339, 347, 351, 358, 366, 368, 370, 372, 374, 376, 378, 383, 388, 396, 399, 408, 411, 415, 423, 430, 439, 452, 455, 458, 461, 464, 467, 470, 476, 479, 485, 488, 494, 498, 501, 505, 510, 515, 521, 526, 530, 535, 543, 551, 557, 566, 577, 584, 588, 595, 602, 610, 614, 618, 622, 629, 636, 644, 650, 655, 660, 664, 672, 677, 682, 686, 689, 697, 701, 703, 708, 710, 715, 721, 727, 733, 739, 744, 748, 755, 761, 766, 772, 775, 782, 787, 791, 796, 800, 810, 815, 823, 829, 836, 843, 849, 856, 862, 866, 869}

func (i Op) String() string {
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
