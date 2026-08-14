; ModuleID = 'testdata/ir/cpp_map_ptr_key/source.cpp'
source_filename = "testdata/ir/cpp_map_ptr_key/source.cpp"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-conda-linux-gnu"

%"struct.std::piecewise_construct_t" = type { i8 }
%struct.V = type { ptr }
%"class.std::map" = type { %"class.std::_Rb_tree" }
%"class.std::_Rb_tree" = type { %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Rb_tree_impl" }
%"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Rb_tree_impl" = type { [8 x i8], %"struct.std::_Rb_tree_header" }
%"struct.std::_Rb_tree_header" = type { %"struct.std::_Rb_tree_node_base", i64 }
%"struct.std::_Rb_tree_node_base" = type { i32, ptr, ptr, ptr }
%"struct.std::_Rb_tree_iterator" = type { ptr }
%"struct.std::pair" = type <{ ptr, i32, [4 x i8] }>
%"struct.std::less" = type { i8 }
%"struct.std::_Rb_tree_const_iterator" = type { ptr }
%"class.std::tuple" = type { %"struct.std::_Tuple_impl" }
%"struct.std::_Tuple_impl" = type { %"struct.std::_Head_base" }
%"struct.std::_Head_base" = type { ptr }
%"class.std::tuple.0" = type { i8 }
%"struct.std::_Rb_tree_node" = type { %"struct.std::_Rb_tree_node_base", %"struct.__gnu_cxx::__aligned_membuf" }
%"struct.__gnu_cxx::__aligned_membuf" = type { [16 x i8] }
%"class.std::allocator" = type { i8 }
%"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Alloc_node" = type { ptr }
%"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node" = type { ptr, ptr }
%"struct.std::pair.1" = type { ptr, ptr }
%"struct.std::_Select1st" = type { i8 }
%"struct.std::_Index_tuple" = type { i8 }
%"struct.std::_Index_tuple.3" = type { i8 }

$_ZN1VC2Ev = comdat any

$_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2Ev = comdat any

$_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEixERS6_ = comdat any

$_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv = comdat any

$_ZStneRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_ = comdat any

$_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv = comdat any

$_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv = comdat any

$_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEppEv = comdat any

$_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2ERKS9_ = comdat any

$_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEED2Ev = comdat any

$_ZNK1V2idEv = comdat any

$_ZNK1V4sameEPKS_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EEC2Ev = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EEC2Ev = comdat any

$_ZNSt20_Rb_tree_key_compareISt4lessIPK1VEEC2Ev = comdat any

$_ZNSt15_Rb_tree_headerC2Ev = comdat any

$_ZNSt15_Rb_tree_header8_M_resetEv = comdat any

$_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv = comdat any

$_ZStneRKSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEES8_ = comdat any

$_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv = comdat any

$_ZNKSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEptEv = comdat any

$_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEppEv = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE5beginEv = comdat any

$_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE3endEv = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv = comdat any

$_ZNKSt18_Rb_tree_node_base11_M_base_ptrEv = comdat any

$_ZNKSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv = comdat any

$_ZNK9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE6_M_ptrEv = comdat any

$_ZNK9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE7_M_addrEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EEC2ERKSB_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EEC2ERKSD_ = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_rootEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyERKSB_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_rootEv = comdat any

$_ZNSt20_Rb_tree_key_compareISt4lessIPK1VEEC2ERKS4_ = comdat any

$_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEED2Ev = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_Alloc_nodeC2ERSB_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_baseRKSB_RT0_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_basePSt13_Rb_tree_nodeIS5_ESF_RT0_ = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_begin_nodeEv = comdat any

$_ZNSt18_Rb_tree_node_base10_S_minimumEPS_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_leftmostEv = comdat any

$_ZNSt18_Rb_tree_node_base10_S_maximumEPS_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_rightmostEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_clone_nodeILb0ENSB_11_Alloc_nodeEEEPSt13_Rb_tree_nodeIS5_ESG_RT0_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_Alloc_nodeclIRKS5_EEPSt13_Rb_tree_nodeIS5_EOT_ = comdat any

$_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_create_nodeIJRKS5_EEEPSt13_Rb_tree_nodeIS5_EDpOT_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_get_nodeEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE17_M_construct_nodeIJRKS5_EEEvPSt13_Rb_tree_nodeIS5_EDpOT_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE21_M_get_Node_allocatorEv = comdat any

$_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEE8allocateEmPKv = comdat any

$_ZN9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE6_M_ptrEv = comdat any

$_ZN9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE7_M_addrEv = comdat any

$_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE11_M_node_ptrEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EED2Ev = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_eraseEPSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EED2Ev = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_drop_nodeEPSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE15_M_destroy_nodeEPSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_put_nodeEPSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEE10deallocateEPS7_m = comdat any

$_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE11lower_boundERS6_ = comdat any

$_ZSteqRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_ = comdat any

$_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE8key_compEv = comdat any

$_ZNKSt4lessIPK1VEclES2_S2_ = comdat any

$_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEdeEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE22_M_emplace_hint_uniqueIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEESt17_Rb_tree_iteratorIS5_ESt23_Rb_tree_const_iteratorIS5_EDpOT_ = comdat any

$_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEC2ERKSt17_Rb_tree_iteratorIS5_E = comdat any

$_ZNSt5tupleIJRKPK1VEEC2ILb1ETnNSt9enable_ifIXclsr4_TCCIXT_EEE29__is_implicitly_constructibleIS4_EEEbE4typeELb1EEES4_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11lower_boundERS4_ = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_lower_boundEPSt18_Rb_tree_node_baseSD_RS4_ = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_beginEv = comdat any

$_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt18_Rb_tree_node_base = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt18_Rb_tree_node_base = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyERKSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNKSt10_Select1stISt4pairIKPK1VjEEclERKS5_ = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8key_compEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_nodeC2IJRKSt21piecewise_construct_tSt5tupleIJRS4_EESH_IJEEEEERSB_DpOT_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE29_M_get_insert_hint_unique_posESt23_Rb_tree_const_iteratorIS5_ERS4_ = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_node6_M_keyEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_node9_M_insertES3_IPSt18_Rb_tree_node_baseSE_E = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_nodeD2Ev = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_create_nodeIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEEPSt13_Rb_tree_nodeIS5_EDpOT_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE17_M_construct_nodeIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEEvPSt13_Rb_tree_nodeIS5_EDpOT_ = comdat any

$_ZNSt5tupleIJRKPK1VEEC2EOS5_ = comdat any

$_ZNSt4pairIKPK1VjEC2IJRS3_EJEEESt21piecewise_construct_tSt5tupleIJDpT_EES8_IJDpT0_EE = comdat any

$_ZNSt11_Tuple_implILm0EJRKPK1VEEC2EOS5_ = comdat any

$_ZNSt4pairIKPK1VjEC2IJRS3_EJLm0EEJETpTnmJEEERSt5tupleIJDpT_EERS7_IJDpT1_EESt12_Index_tupleIJXspT0_EEESG_IJXspT2_EEE = comdat any

$_ZSt3getILm0EJRKPK1VEERNSt13tuple_elementIXT_ESt5tupleIJDpT0_EEE4typeERS9_ = comdat any

$_ZSt12__get_helperILm0ERKPK1VJEERT0_RSt11_Tuple_implIXT_EJS5_DpT1_EE = comdat any

$_ZNSt11_Tuple_implILm0EJRKPK1VEE7_M_headERS5_ = comdat any

$_ZNSt10_Head_baseILm0ERKPK1VLb0EE7_M_headERS5_ = comdat any

$_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE4sizeEv = comdat any

$_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IS1_RS1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE24_M_get_insert_unique_posERS4_ = comdat any

$_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S4_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_ = comdat any

$_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEmmEv = comdat any

$_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE5beginEv = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_insert_nodeEPSt18_Rb_tree_node_baseSD_PSt13_Rb_tree_nodeIS5_E = comdat any

$_ZNSt9__rb_tree12_Node_traitsISt4pairIKPK1VjEPS6_E23_S_insert_and_rebalanceEbPSt18_Rb_tree_node_baseSA_RS9_ = comdat any

$_ZNSt11_Tuple_implILm0EJRKPK1VEEC2ES4_ = comdat any

$_ZNSt10_Head_baseILm0ERKPK1VLb0EEC2ES4_ = comdat any

$_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE3endEv = comdat any

$_ZTV1V = comdat any

$_ZSt19piecewise_construct = comdat any

@.str = private unnamed_addr constant [9 x i8] c"nil key\0A\00", align 1
@.str.1 = private unnamed_addr constant [7 x i8] c"%u %u\0A\00", align 1
@.str.2 = private unnamed_addr constant [6 x i8] c"miss\0A\00", align 1
@.str.3 = private unnamed_addr constant [14 x i8] c"nil key copy\0A\00", align 1
@.str.4 = private unnamed_addr constant [16 x i8] c"nil key memcpy\0A\00", align 1
@.str.5 = private unnamed_addr constant [4 x i8] c"%u\0A\00", align 1
@_ZTV1V = linkonce_odr dso_local unnamed_addr constant { [4 x ptr] } { [4 x ptr] [ptr null, ptr null, ptr @_ZNK1V2idEv, ptr @_ZNK1V4sameEPKS_] }, comdat, align 8
@_ZSt19piecewise_construct = linkonce_odr dso_local constant %"struct.std::piecewise_construct_t" zeroinitializer, comdat, align 1

; Function Attrs: mustprogress noinline norecurse nounwind optnone uwtable
define dso_local noundef i32 @main() #0 {
entry:
  %retval = alloca i32, align 4
  %a = alloca %struct.V, align 8
  %b = alloca %struct.V, align 8
  %pa = alloca ptr, align 8
  %pb = alloca ptr, align 8
  %m = alloca %"class.std::map", align 8
  %n = alloca i32, align 4
  %sum = alloca i32, align 4
  %it = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp = alloca %"struct.std::_Rb_tree_iterator", align 8
  %cleanup.dest.slot = alloca i32, align 4
  %copied = alloca %"class.std::map", align 8
  %cn = alloca i32, align 4
  %csum = alloca i32, align 4
  %it20 = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp24 = alloca %"struct.std::_Rb_tree_iterator", align 8
  %empty_src = alloca %"class.std::map", align 8
  %empty_dst = alloca %"class.std::map", align 8
  %en = alloca i32, align 4
  %it49 = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp53 = alloca %"struct.std::_Rb_tree_iterator", align 8
  store i32 0, ptr %retval, align 4
  call void @_ZN1VC2Ev(ptr noundef nonnull align 8 dereferenceable(8) %a) #10
  call void @_ZN1VC2Ev(ptr noundef nonnull align 8 dereferenceable(8) %b) #10
  store ptr %a, ptr %pa, align 8
  store ptr %b, ptr %pb, align 8
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %m) #10
  %call = call noundef nonnull align 4 dereferenceable(4) ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEixERS6_(ptr noundef nonnull align 8 dereferenceable(48) %m, ptr noundef nonnull align 8 dereferenceable(8) %pa)
  store i32 1, ptr %call, align 4
  %call1 = call noundef nonnull align 4 dereferenceable(4) ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEixERS6_(ptr noundef nonnull align 8 dereferenceable(48) %m, ptr noundef nonnull align 8 dereferenceable(8) %pb)
  store i32 2, ptr %call1, align 4
  store i32 0, ptr %n, align 4
  store i32 0, ptr %sum, align 4
  %call2 = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %m) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %it, i32 0, i32 0
  store ptr %call2, ptr %coerce.dive, align 8
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %call3 = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %m) #10
  %coerce.dive4 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %ref.tmp, i32 0, i32 0
  store ptr %call3, ptr %coerce.dive4, align 8
  %call5 = call noundef zeroext i1 @_ZStneRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %it, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp) #10
  br i1 %call5, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %call6 = call noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it) #10
  %first = getelementptr inbounds nuw %"struct.std::pair", ptr %call6, i32 0, i32 0
  %0 = load ptr, ptr %first, align 8
  %cmp = icmp eq ptr %0, null
  br i1 %cmp, label %if.then, label %if.end

if.then:                                          ; preds = %for.body
  %call7 = call i32 (ptr, ...) @printf(ptr noundef @.str)
  store i32 1, ptr %retval, align 4
  store i32 1, ptr %cleanup.dest.slot, align 4
  br label %cleanup71

if.end:                                           ; preds = %for.body
  %call8 = call noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it) #10
  %second = getelementptr inbounds nuw %"struct.std::pair", ptr %call8, i32 0, i32 1
  %1 = load i32, ptr %second, align 8
  %call9 = call noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it) #10
  %first10 = getelementptr inbounds nuw %"struct.std::pair", ptr %call9, i32 0, i32 0
  %2 = load ptr, ptr %first10, align 8
  %vtable = load ptr, ptr %2, align 8
  %vfn = getelementptr inbounds ptr, ptr %vtable, i64 0
  %3 = load ptr, ptr %vfn, align 8
  %call11 = call noundef i32 %3(ptr noundef nonnull align 8 dereferenceable(8) %2)
  %add = add i32 %1, %call11
  %4 = load i32, ptr %sum, align 4
  %add12 = add i32 %4, %add
  store i32 %add12, ptr %sum, align 4
  %5 = load i32, ptr %n, align 4
  %inc = add i32 %5, 1
  store i32 %inc, ptr %n, align 4
  br label %for.inc

for.inc:                                          ; preds = %if.end
  %call13 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEppEv(ptr noundef nonnull align 8 dereferenceable(8) %it) #10
  br label %for.cond, !llvm.loop !6

for.end:                                          ; preds = %for.cond
  %6 = load i32, ptr %n, align 4
  %7 = load i32, ptr %sum, align 4
  %call14 = call i32 (ptr, ...) @printf(ptr noundef @.str.1, i32 noundef %6, i32 noundef %7)
  %8 = load ptr, ptr %pa, align 8
  %call15 = call noundef zeroext i1 @_ZL7touchesRKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEES2_(ptr noundef nonnull align 8 dereferenceable(48) %m, ptr noundef %8)
  br i1 %call15, label %lor.lhs.false, label %if.then17

lor.lhs.false:                                    ; preds = %for.end
  %9 = load ptr, ptr %pb, align 8
  %call16 = call noundef zeroext i1 @_ZL7touchesRKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEES2_(ptr noundef nonnull align 8 dereferenceable(48) %m, ptr noundef %9)
  br i1 %call16, label %if.end19, label %if.then17

if.then17:                                        ; preds = %lor.lhs.false, %for.end
  %call18 = call i32 (ptr, ...) @printf(ptr noundef @.str.2)
  store i32 1, ptr %retval, align 4
  store i32 1, ptr %cleanup.dest.slot, align 4
  br label %cleanup71

if.end19:                                         ; preds = %lor.lhs.false
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2ERKS9_(ptr noundef nonnull align 8 dereferenceable(48) %copied, ptr noundef nonnull align 8 dereferenceable(48) %m)
  store i32 0, ptr %cn, align 4
  store i32 0, ptr %csum, align 4
  %call21 = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %copied) #10
  %coerce.dive22 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %it20, i32 0, i32 0
  store ptr %call21, ptr %coerce.dive22, align 8
  br label %for.cond23

for.cond23:                                       ; preds = %for.inc45, %if.end19
  %call25 = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %copied) #10
  %coerce.dive26 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %ref.tmp24, i32 0, i32 0
  store ptr %call25, ptr %coerce.dive26, align 8
  %call27 = call noundef zeroext i1 @_ZStneRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %it20, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp24) #10
  br i1 %call27, label %for.body28, label %for.end47

for.body28:                                       ; preds = %for.cond23
  %call29 = call noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it20) #10
  %first30 = getelementptr inbounds nuw %"struct.std::pair", ptr %call29, i32 0, i32 0
  %10 = load ptr, ptr %first30, align 8
  %cmp31 = icmp eq ptr %10, null
  br i1 %cmp31, label %if.then32, label %if.end34

if.then32:                                        ; preds = %for.body28
  %call33 = call i32 (ptr, ...) @printf(ptr noundef @.str.3)
  store i32 1, ptr %retval, align 4
  store i32 1, ptr %cleanup.dest.slot, align 4
  br label %cleanup70

if.end34:                                         ; preds = %for.body28
  %call35 = call noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it20) #10
  %second36 = getelementptr inbounds nuw %"struct.std::pair", ptr %call35, i32 0, i32 1
  %11 = load i32, ptr %second36, align 8
  %call37 = call noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it20) #10
  %first38 = getelementptr inbounds nuw %"struct.std::pair", ptr %call37, i32 0, i32 0
  %12 = load ptr, ptr %first38, align 8
  %vtable39 = load ptr, ptr %12, align 8
  %vfn40 = getelementptr inbounds ptr, ptr %vtable39, i64 0
  %13 = load ptr, ptr %vfn40, align 8
  %call41 = call noundef i32 %13(ptr noundef nonnull align 8 dereferenceable(8) %12)
  %add42 = add i32 %11, %call41
  %14 = load i32, ptr %csum, align 4
  %add43 = add i32 %14, %add42
  store i32 %add43, ptr %csum, align 4
  %15 = load i32, ptr %cn, align 4
  %inc44 = add i32 %15, 1
  store i32 %inc44, ptr %cn, align 4
  br label %for.inc45

for.inc45:                                        ; preds = %if.end34
  %call46 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEppEv(ptr noundef nonnull align 8 dereferenceable(8) %it20) #10
  br label %for.cond23, !llvm.loop !8

for.end47:                                        ; preds = %for.cond23
  %16 = load i32, ptr %cn, align 4
  %17 = load i32, ptr %csum, align 4
  %call48 = call i32 (ptr, ...) @printf(ptr noundef @.str.1, i32 noundef %16, i32 noundef %17)
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %empty_src) #10
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %empty_dst) #10
  call void @llvm.memcpy.p0.p0.i64(ptr align 1 %empty_dst, ptr align 1 %empty_src, i64 48, i1 false)
  store i32 0, ptr %en, align 4
  %call50 = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %empty_dst) #10
  %coerce.dive51 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %it49, i32 0, i32 0
  store ptr %call50, ptr %coerce.dive51, align 8
  br label %for.cond52

for.cond52:                                       ; preds = %for.inc65, %for.end47
  %call54 = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %empty_dst) #10
  %coerce.dive55 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %ref.tmp53, i32 0, i32 0
  store ptr %call54, ptr %coerce.dive55, align 8
  %call56 = call noundef zeroext i1 @_ZStneRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %it49, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp53) #10
  br i1 %call56, label %for.body57, label %for.end67

for.body57:                                       ; preds = %for.cond52
  %call58 = call noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it49) #10
  %first59 = getelementptr inbounds nuw %"struct.std::pair", ptr %call58, i32 0, i32 0
  %18 = load ptr, ptr %first59, align 8
  %cmp60 = icmp eq ptr %18, null
  br i1 %cmp60, label %if.then61, label %if.end63

if.then61:                                        ; preds = %for.body57
  %call62 = call i32 (ptr, ...) @printf(ptr noundef @.str.4)
  store i32 1, ptr %retval, align 4
  store i32 1, ptr %cleanup.dest.slot, align 4
  br label %cleanup

if.end63:                                         ; preds = %for.body57
  %19 = load i32, ptr %en, align 4
  %inc64 = add i32 %19, 1
  store i32 %inc64, ptr %en, align 4
  br label %for.inc65

for.inc65:                                        ; preds = %if.end63
  %call66 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEppEv(ptr noundef nonnull align 8 dereferenceable(8) %it49) #10
  br label %for.cond52, !llvm.loop !9

for.end67:                                        ; preds = %for.cond52
  %20 = load i32, ptr %en, align 4
  %call68 = call i32 (ptr, ...) @printf(ptr noundef @.str.5, i32 noundef %20)
  store i32 0, ptr %retval, align 4
  store i32 1, ptr %cleanup.dest.slot, align 4
  br label %cleanup

cleanup:                                          ; preds = %for.end67, %if.then61
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %empty_dst) #10
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %empty_src) #10
  br label %cleanup70

cleanup70:                                        ; preds = %cleanup, %if.then32
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %copied) #10
  br label %cleanup71

cleanup71:                                        ; preds = %cleanup70, %if.then17, %if.then
  call void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %m) #10
  %21 = load i32, ptr %retval, align 4
  ret i32 %21
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZN1VC2Ev(ptr noundef nonnull align 8 dereferenceable(8) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  store ptr getelementptr inbounds inrange(-16, 16) ({ [4 x ptr] }, ptr @_ZTV1V, i32 0, i32 0, i32 2), ptr %this1, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %_M_t) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 4 dereferenceable(4) ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEixERS6_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(8) %__k) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__k.addr = alloca ptr, align 8
  %__i = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp5 = alloca %"struct.std::less", align 1
  %undef.agg.tmp = alloca %"struct.std::less", align 1
  %ref.tmp8 = alloca %"struct.std::_Rb_tree_iterator", align 8
  %agg.tmp = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %ref.tmp9 = alloca %"class.std::tuple", align 8
  %ref.tmp10 = alloca %"class.std::tuple.0", align 1
  store ptr %this, ptr %this.addr, align 8
  store ptr %__k, ptr %__k.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE11lower_boundERS6_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %0)
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %__i, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %call2 = call ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %coerce.dive3 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %ref.tmp, i32 0, i32 0
  store ptr %call2, ptr %coerce.dive3, align 8
  %call4 = call noundef zeroext i1 @_ZSteqRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %__i, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp) #10
  br i1 %call4, label %lor.end, label %lor.rhs

lor.rhs:                                          ; preds = %entry
  call void @_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE8key_compEv(ptr noundef nonnull align 8 dereferenceable(48) %this1)
  %1 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %2 = load ptr, ptr %1, align 8
  %call6 = call noundef nonnull align 8 dereferenceable(12) ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEdeEv(ptr noundef nonnull align 8 dereferenceable(8) %__i) #10
  %first = getelementptr inbounds nuw %"struct.std::pair", ptr %call6, i32 0, i32 0
  %3 = load ptr, ptr %first, align 8
  %call7 = call noundef zeroext i1 @_ZNKSt4lessIPK1VEclES2_S2_(ptr noundef nonnull align 1 dereferenceable(1) %ref.tmp5, ptr noundef %2, ptr noundef %3) #10
  br label %lor.end

lor.end:                                          ; preds = %lor.rhs, %entry
  %4 = phi i1 [ true, %entry ], [ %call7, %lor.rhs ]
  br i1 %4, label %if.then, label %if.end

if.then:                                          ; preds = %lor.end
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  call void @_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEC2ERKSt17_Rb_tree_iteratorIS5_E(ptr noundef nonnull align 8 dereferenceable(8) %agg.tmp, ptr noundef nonnull align 8 dereferenceable(8) %__i) #10
  %5 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  call void @_ZNSt5tupleIJRKPK1VEEC2ILb1ETnNSt9enable_ifIXclsr4_TCCIXT_EEE29__is_implicitly_constructibleIS4_EEEbE4typeELb1EEES4_(ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp9, ptr noundef nonnull align 8 dereferenceable(8) %5) #10
  %coerce.dive11 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %agg.tmp, i32 0, i32 0
  %6 = load ptr, ptr %coerce.dive11, align 8
  %call12 = call ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE22_M_emplace_hint_uniqueIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEESt17_Rb_tree_iteratorIS5_ESt23_Rb_tree_const_iteratorIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %_M_t, ptr %6, ptr noundef nonnull align 1 dereferenceable(1) @_ZSt19piecewise_construct, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp9, ptr noundef nonnull align 1 dereferenceable(1) %ref.tmp10)
  %coerce.dive13 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %ref.tmp8, i32 0, i32 0
  store ptr %call12, ptr %coerce.dive13, align 8
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %__i, ptr align 8 %ref.tmp8, i64 8, i1 false)
  br label %if.end

if.end:                                           ; preds = %if.then, %lor.end
  %call14 = call noundef nonnull align 8 dereferenceable(12) ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEdeEv(ptr noundef nonnull align 8 dereferenceable(8) %__i) #10
  %second = getelementptr inbounds nuw %"struct.std::pair", ptr %call14, i32 0, i32 1
  ret ptr %second
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  %call = call ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %_M_t) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %coerce.dive2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %0 = load ptr, ptr %coerce.dive2, align 8
  ret ptr %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef zeroext i1 @_ZStneRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y) #1 comdat {
entry:
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %0, i32 0, i32 0
  %1 = load ptr, ptr %_M_node, align 8
  %2 = load ptr, ptr %__y.addr, align 8, !nonnull !10, !align !11
  %_M_node1 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %2, i32 0, i32 0
  %3 = load ptr, ptr %_M_node1, align 8
  %cmp = icmp ne ptr %1, %3
  ret i1 %cmp
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  %call = call ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %_M_t) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %coerce.dive2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %0 = load ptr, ptr %coerce.dive2, align 8
  ret ptr %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %0)
  ret ptr %call
}

declare i32 @printf(ptr noundef, ...) #2

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEppEv(ptr noundef nonnull align 8 dereferenceable(8) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef ptr @_ZSt18_Rb_tree_incrementPSt18_Rb_tree_node_base(ptr noundef %0) #11
  %_M_node2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %this1, i32 0, i32 0
  store ptr %call, ptr %_M_node2, align 8
  ret ptr %this1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define internal noundef zeroext i1 @_ZL7touchesRKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEES2_(ptr noundef nonnull align 8 dereferenceable(48) %m, ptr noundef %v) #1 {
entry:
  %retval = alloca i1, align 1
  %m.addr = alloca ptr, align 8
  %v.addr = alloca ptr, align 8
  %it = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %ref.tmp = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  store ptr %m, ptr %m.addr, align 8
  store ptr %v, ptr %v.addr, align 8
  %0 = load ptr, ptr %m.addr, align 8, !nonnull !10, !align !11
  %call = call ptr @_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %0) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %it, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  br label %for.cond

for.cond:                                         ; preds = %for.inc, %entry
  %1 = load ptr, ptr %m.addr, align 8, !nonnull !10, !align !11
  %call1 = call ptr @_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %1) #10
  %coerce.dive2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %ref.tmp, i32 0, i32 0
  store ptr %call1, ptr %coerce.dive2, align 8
  %call3 = call noundef zeroext i1 @_ZStneRKSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %it, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp) #10
  br i1 %call3, label %for.body, label %for.end

for.body:                                         ; preds = %for.cond
  %2 = load ptr, ptr %v.addr, align 8
  %call4 = call noundef ptr @_ZNKSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %it) #10
  %first = getelementptr inbounds nuw %"struct.std::pair", ptr %call4, i32 0, i32 0
  %3 = load ptr, ptr %first, align 8
  %vtable = load ptr, ptr %2, align 8
  %vfn = getelementptr inbounds ptr, ptr %vtable, i64 1
  %4 = load ptr, ptr %vfn, align 8
  %call5 = call noundef zeroext i1 %4(ptr noundef nonnull align 8 dereferenceable(8) %2, ptr noundef %3)
  br i1 %call5, label %if.then, label %if.end

if.then:                                          ; preds = %for.body
  store i1 true, ptr %retval, align 1
  br label %return

if.end:                                           ; preds = %for.body
  br label %for.inc

for.inc:                                          ; preds = %if.end
  %call6 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEppEv(ptr noundef nonnull align 8 dereferenceable(8) %it) #10
  br label %for.cond, !llvm.loop !12

for.end:                                          ; preds = %for.cond
  store i1 false, ptr %retval, align 1
  br label %return

return:                                           ; preds = %for.end, %if.then
  %5 = load i1, ptr %retval, align 1
  ret i1 %5
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEEC2ERKS9_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(48) %0) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %0, ptr %.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  %1 = load ptr, ptr %.addr, align 8, !nonnull !10, !align !11
  %_M_t2 = getelementptr inbounds nuw %"class.std::map", ptr %1, i32 0, i32 0
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EEC2ERKSB_(ptr noundef nonnull align 8 dereferenceable(48) %_M_t, ptr noundef nonnull align 8 dereferenceable(48) %_M_t2)
  ret void
}

; Function Attrs: nocallback nofree nounwind willreturn memory(argmem: readwrite)
declare void @llvm.memcpy.p0.p0.i64(ptr noalias writeonly captures(none), ptr noalias readonly captures(none), i64, i1 immarg) #3

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %_M_t) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef i32 @_ZNK1V2idEv(ptr noundef nonnull align 8 dereferenceable(8) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  ret i32 42
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef zeroext i1 @_ZNK1V4sameEPKS_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef %o) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %o.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %o, ptr %o.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %o.addr, align 8
  %cmp = icmp eq ptr %this1, %0
  ret i1 %cmp
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %_M_impl) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EEC2Ev(ptr noundef nonnull align 8 dereferenceable(48) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr.i2 = alloca ptr, align 8
  %this.addr.i = alloca ptr, align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  store ptr %this1, ptr %this.addr.i, align 8
  %this1.i = load ptr, ptr %this.addr.i, align 8
  store ptr %this1.i, ptr %this.addr.i2, align 8
  %this1.i3 = load ptr, ptr %this.addr.i2, align 8
  call void @_ZNSt20_Rb_tree_key_compareISt4lessIPK1VEEC2Ev(ptr noundef nonnull align 1 dereferenceable(1) %this1) #10
  %0 = getelementptr inbounds i8, ptr %this1, i64 8
  call void @_ZNSt15_Rb_tree_headerC2Ev(ptr noundef nonnull align 8 dereferenceable(40) %0) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt20_Rb_tree_key_compareISt4lessIPK1VEEC2Ev(ptr noundef nonnull align 1 dereferenceable(1) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt15_Rb_tree_headerC2Ev(ptr noundef nonnull align 8 dereferenceable(40) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 0
  %_M_header2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 0
  %_M_color = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header2, i32 0, i32 0
  store i32 0, ptr %_M_color, align 8
  call void @_ZNSt15_Rb_tree_header8_M_resetEv(ptr noundef nonnull align 8 dereferenceable(40) %this1)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt15_Rb_tree_header8_M_resetEv(ptr noundef nonnull align 8 dereferenceable(40) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 0
  %_M_parent = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 1
  store ptr null, ptr %_M_parent, align 8
  %_M_header2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 0
  %_M_header3 = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 0
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header3, i32 0, i32 2
  store ptr %_M_header2, ptr %_M_left, align 8
  %_M_header4 = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 0
  %_M_header5 = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 0
  %_M_right = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header5, i32 0, i32 3
  store ptr %_M_header4, ptr %_M_right, align 8
  %_M_node_count = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %this1, i32 0, i32 1
  store i64 0, ptr %_M_node_count, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  %call = call ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %_M_t) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %retval, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %coerce.dive2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %retval, i32 0, i32 0
  %0 = load ptr, ptr %coerce.dive2, align 8
  ret ptr %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef zeroext i1 @_ZStneRKSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y) #1 comdat {
entry:
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %0, i32 0, i32 0
  %1 = load ptr, ptr %_M_node, align 8
  %2 = load ptr, ptr %__y.addr, align 8, !nonnull !10, !align !11
  %_M_node1 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %2, i32 0, i32 0
  %3 = load ptr, ptr %_M_node1, align 8
  %cmp = icmp ne ptr %1, %3
  ret i1 %cmp
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  %call = call ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %_M_t) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %retval, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %coerce.dive2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %retval, i32 0, i32 0
  %0 = load ptr, ptr %coerce.dive2, align 8
  ret ptr %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEptEv(ptr noundef nonnull align 8 dereferenceable(8) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef ptr @_ZNKSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %0)
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEppEv(ptr noundef nonnull align 8 dereferenceable(8) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef ptr @_ZSt18_Rb_tree_incrementPSt18_Rb_tree_node_base(ptr noundef %0) #11
  %_M_node2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %this1, i32 0, i32 0
  store ptr %call, ptr %_M_node2, align 8
  ret ptr %this1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 2
  %0 = load ptr, ptr %_M_left, align 8
  call void @_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %retval, ptr noundef %0) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %retval, i32 0, i32 0
  %1 = load ptr, ptr %coerce.dive, align 8
  ret ptr %1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef %__x) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__x.addr, align 8
  store ptr %0, ptr %_M_node, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  call void @_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %retval, ptr noundef %call) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %retval, i32 0, i32 0
  %0 = load ptr, ptr %coerce.dive, align 8
  ret ptr %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %call = call noundef ptr @_ZNKSt18_Rb_tree_node_base11_M_base_ptrEv(ptr noundef nonnull align 8 dereferenceable(32) %_M_header) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt18_Rb_tree_node_base11_M_base_ptrEv(ptr noundef nonnull align 8 dereferenceable(32) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  ret ptr %this1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_storage = getelementptr inbounds nuw %"struct.std::_Rb_tree_node", ptr %this1, i32 0, i32 1
  %call = call noundef ptr @_ZNK9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE6_M_ptrEv(ptr noundef nonnull align 8 dereferenceable(16) %_M_storage) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNK9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE6_M_ptrEv(ptr noundef nonnull align 8 dereferenceable(16) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNK9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE7_M_addrEv(ptr noundef nonnull align 8 dereferenceable(16) %this1) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNK9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE7_M_addrEv(ptr noundef nonnull align 8 dereferenceable(16) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_storage = getelementptr inbounds nuw %"struct.__gnu_cxx::__aligned_membuf", ptr %this1, i32 0, i32 0
  ret ptr %_M_storage
}

; Function Attrs: nounwind willreturn memory(read)
declare noundef ptr @_ZSt18_Rb_tree_incrementPSt18_Rb_tree_node_base(ptr noundef) #4

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EEC2ERKSB_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(48) %__x) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %_M_impl2 = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %0, i32 0, i32 0
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EEC2ERKSD_(ptr noundef nonnull align 8 dereferenceable(48) %_M_impl, ptr noundef nonnull align 8 dereferenceable(48) %_M_impl2)
  %1 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_rootEv(ptr noundef nonnull align 8 dereferenceable(48) %1) #10
  %tobool = icmp ne ptr %call, null
  br i1 %tobool, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %2 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %call3 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyERKSB_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(48) %2)
  %call4 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_rootEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  store ptr %call3, ptr %call4, align 8
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EEC2ERKSD_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(48) %__x) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr.i11 = alloca ptr, align 8
  %.addr.i12 = alloca ptr, align 8
  %this.addr.i9 = alloca ptr, align 8
  %.addr.i = alloca ptr, align 8
  %this.addr.i6 = alloca ptr, align 8
  %__a.addr.i7 = alloca ptr, align 8
  %this.addr.i3 = alloca ptr, align 8
  %__a.addr.i4 = alloca ptr, align 8
  %result.ptr.i2 = alloca ptr, align 8
  %__rhs.addr.i = alloca ptr, align 8
  %this.addr.i = alloca ptr, align 8
  %result.ptr.i = alloca ptr, align 8
  %__a.addr.i = alloca ptr, align 8
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %ref.tmp = alloca %"class.std::allocator", align 1
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  call void @llvm.experimental.noalias.scope.decl(metadata !13)
  store ptr %ref.tmp, ptr %result.ptr.i, align 8, !noalias !13
  store ptr %0, ptr %__a.addr.i, align 8, !noalias !13
  %1 = load ptr, ptr %__a.addr.i, align 8, !noalias !13, !nonnull !10
  call void @llvm.experimental.noalias.scope.decl(metadata !16)
  store ptr %ref.tmp, ptr %result.ptr.i2, align 8, !noalias !16
  store ptr %1, ptr %__rhs.addr.i, align 8, !noalias !16
  %2 = load ptr, ptr %__rhs.addr.i, align 8, !noalias !16, !nonnull !10
  store ptr %ref.tmp, ptr %this.addr.i3, align 8
  store ptr %2, ptr %__a.addr.i4, align 8
  %this1.i5 = load ptr, ptr %this.addr.i3, align 8
  %3 = load ptr, ptr %__a.addr.i4, align 8, !nonnull !10
  store ptr %this1.i5, ptr %this.addr.i11, align 8
  store ptr %3, ptr %.addr.i12, align 8
  %this1.i13 = load ptr, ptr %this.addr.i11, align 8
  store ptr %this1, ptr %this.addr.i6, align 8
  store ptr %ref.tmp, ptr %__a.addr.i7, align 8
  %this1.i8 = load ptr, ptr %this.addr.i6, align 8
  %4 = load ptr, ptr %__a.addr.i7, align 8, !nonnull !10
  store ptr %this1.i8, ptr %this.addr.i9, align 8
  store ptr %4, ptr %.addr.i, align 8
  %this1.i10 = load ptr, ptr %this.addr.i9, align 8
  store ptr %ref.tmp, ptr %this.addr.i, align 8
  %this1.i = load ptr, ptr %this.addr.i, align 8
  call void @_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEED2Ev(ptr noundef nonnull align 1 dereferenceable(1) %this1.i) #10
  %5 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  call void @_ZNSt20_Rb_tree_key_compareISt4lessIPK1VEEC2ERKS4_(ptr noundef nonnull align 1 dereferenceable(1) %this1, ptr noundef nonnull align 1 dereferenceable(1) %5)
  %6 = getelementptr inbounds i8, ptr %this1, i64 8
  call void @_ZNSt15_Rb_tree_headerC2Ev(ptr noundef nonnull align 8 dereferenceable(40) %6) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_rootEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_parent = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 1
  %0 = load ptr, ptr %_M_parent, align 8
  ret ptr %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyERKSB_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(48) %__x) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__an = alloca %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Alloc_node", align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_Alloc_nodeC2ERSB_(ptr noundef nonnull align 8 dereferenceable(8) %__an, ptr noundef nonnull align 8 dereferenceable(48) %this1)
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %call = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_baseRKSB_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(48) %0, ptr noundef nonnull align 8 dereferenceable(8) %__an)
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_rootEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_parent = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 1
  ret ptr %_M_parent
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt20_Rb_tree_key_compareISt4lessIPK1VEEC2ERKS4_(ptr noundef nonnull align 1 dereferenceable(1) %this, ptr noundef nonnull align 1 dereferenceable(1) %__comp) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__comp.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__comp, ptr %__comp.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__comp.addr, align 8, !nonnull !10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEED2Ev(ptr noundef nonnull align 1 dereferenceable(1) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_Alloc_nodeC2ERSB_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(48) %__t) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__t.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__t, ptr %__t.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Alloc_node", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__t.addr, align 8, !nonnull !10, !align !11
  store ptr %0, ptr %_M_t, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_baseRKSB_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(48) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__gen) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__gen.addr = alloca ptr, align 8
  %__root = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__gen, ptr %__gen.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_begin_nodeEv(ptr noundef nonnull align 8 dereferenceable(48) %0) #10
  %call2 = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %1 = load ptr, ptr %__gen.addr, align 8, !nonnull !10, !align !11
  %call3 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_basePSt13_Rb_tree_nodeIS5_ESF_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %call, ptr noundef %call2, ptr noundef nonnull align 8 dereferenceable(8) %1)
  store ptr %call3, ptr %__root, align 8
  %2 = load ptr, ptr %__root, align 8
  %call4 = call noundef ptr @_ZNSt18_Rb_tree_node_base10_S_minimumEPS_(ptr noundef %2) #10
  %call5 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_leftmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  store ptr %call4, ptr %call5, align 8
  %3 = load ptr, ptr %__root, align 8
  %call6 = call noundef ptr @_ZNSt18_Rb_tree_node_base10_S_maximumEPS_(ptr noundef %3) #10
  %call7 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_rightmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  store ptr %call6, ptr %call7, align 8
  %4 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %4, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_node_count = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 1
  %5 = load i64, ptr %_M_node_count, align 8
  %_M_impl8 = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr9 = getelementptr inbounds i8, ptr %_M_impl8, i64 8
  %_M_node_count10 = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr9, i32 0, i32 1
  store i64 %5, ptr %_M_node_count10, align 8
  %6 = load ptr, ptr %__root, align 8
  ret ptr %6
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_basePSt13_Rb_tree_nodeIS5_ESF_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__x, ptr noundef %__p, ptr noundef nonnull align 8 dereferenceable(8) %__node_gen) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__p.addr = alloca ptr, align 8
  %__node_gen.addr = alloca ptr, align 8
  %__top = alloca ptr, align 8
  %__top_base = alloca ptr, align 8
  %__y = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__p, ptr %__p.addr, align 8
  store ptr %__node_gen, ptr %__node_gen.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %1 = load ptr, ptr %__node_gen.addr, align 8, !nonnull !10, !align !11
  %call = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_clone_nodeILb0ENSB_11_Alloc_nodeEEEPSt13_Rb_tree_nodeIS5_ESG_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %0, ptr noundef nonnull align 8 dereferenceable(8) %1)
  store ptr %call, ptr %__top, align 8
  %2 = load ptr, ptr %__top, align 8
  %call2 = call noundef ptr @_ZNKSt18_Rb_tree_node_base11_M_base_ptrEv(ptr noundef nonnull align 8 dereferenceable(32) %2) #10
  store ptr %call2, ptr %__top_base, align 8
  %3 = load ptr, ptr %__p.addr, align 8
  %4 = load ptr, ptr %__top, align 8
  %_M_parent = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %4, i32 0, i32 1
  store ptr %3, ptr %_M_parent, align 8
  %5 = load ptr, ptr %__x.addr, align 8
  %_M_right = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %5, i32 0, i32 3
  %6 = load ptr, ptr %_M_right, align 8
  %tobool = icmp ne ptr %6, null
  br i1 %tobool, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %7 = load ptr, ptr %__x.addr, align 8
  %call3 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt13_Rb_tree_nodeIS5_E(ptr noundef %7) #10
  %8 = load ptr, ptr %__top_base, align 8
  %9 = load ptr, ptr %__node_gen.addr, align 8, !nonnull !10, !align !11
  %call4 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_basePSt13_Rb_tree_nodeIS5_ESF_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %call3, ptr noundef %8, ptr noundef nonnull align 8 dereferenceable(8) %9)
  %10 = load ptr, ptr %__top, align 8
  %_M_right5 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %10, i32 0, i32 3
  store ptr %call4, ptr %_M_right5, align 8
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  %11 = load ptr, ptr %__top_base, align 8
  store ptr %11, ptr %__p.addr, align 8
  %12 = load ptr, ptr %__x.addr, align 8
  %call6 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt13_Rb_tree_nodeIS5_E(ptr noundef %12)
  store ptr %call6, ptr %__x.addr, align 8
  br label %while.cond

while.cond:                                       ; preds = %if.end17, %if.end
  %13 = load ptr, ptr %__x.addr, align 8
  %tobool7 = icmp ne ptr %13, null
  br i1 %tobool7, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %14 = load ptr, ptr %__x.addr, align 8
  %15 = load ptr, ptr %__node_gen.addr, align 8, !nonnull !10, !align !11
  %call8 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_clone_nodeILb0ENSB_11_Alloc_nodeEEEPSt13_Rb_tree_nodeIS5_ESG_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %14, ptr noundef nonnull align 8 dereferenceable(8) %15)
  %call9 = call noundef ptr @_ZNKSt18_Rb_tree_node_base11_M_base_ptrEv(ptr noundef nonnull align 8 dereferenceable(32) %call8) #10
  store ptr %call9, ptr %__y, align 8
  %16 = load ptr, ptr %__y, align 8
  %17 = load ptr, ptr %__p.addr, align 8
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %17, i32 0, i32 2
  store ptr %16, ptr %_M_left, align 8
  %18 = load ptr, ptr %__p.addr, align 8
  %19 = load ptr, ptr %__y, align 8
  %_M_parent10 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %19, i32 0, i32 1
  store ptr %18, ptr %_M_parent10, align 8
  %20 = load ptr, ptr %__x.addr, align 8
  %_M_right11 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %20, i32 0, i32 3
  %21 = load ptr, ptr %_M_right11, align 8
  %tobool12 = icmp ne ptr %21, null
  br i1 %tobool12, label %if.then13, label %if.end17

if.then13:                                        ; preds = %while.body
  %22 = load ptr, ptr %__x.addr, align 8
  %call14 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt13_Rb_tree_nodeIS5_E(ptr noundef %22) #10
  %23 = load ptr, ptr %__y, align 8
  %24 = load ptr, ptr %__node_gen.addr, align 8, !nonnull !10, !align !11
  %call15 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_M_copyILb0ENSB_11_Alloc_nodeEEEPSt18_Rb_tree_node_basePSt13_Rb_tree_nodeIS5_ESF_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %call14, ptr noundef %23, ptr noundef nonnull align 8 dereferenceable(8) %24)
  %25 = load ptr, ptr %__y, align 8
  %_M_right16 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %25, i32 0, i32 3
  store ptr %call15, ptr %_M_right16, align 8
  br label %if.end17

if.end17:                                         ; preds = %if.then13, %while.body
  %26 = load ptr, ptr %__y, align 8
  store ptr %26, ptr %__p.addr, align 8
  %27 = load ptr, ptr %__x.addr, align 8
  %call18 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt13_Rb_tree_nodeIS5_E(ptr noundef %27)
  store ptr %call18, ptr %__x.addr, align 8
  br label %while.cond, !llvm.loop !19

while.end:                                        ; preds = %while.cond
  %28 = load ptr, ptr %__top_base, align 8
  ret ptr %28
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_begin_nodeEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__begin = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_parent = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 1
  %0 = load ptr, ptr %_M_parent, align 8
  store ptr %0, ptr %__begin, align 8
  %1 = load ptr, ptr %__begin, align 8
  %tobool = icmp ne ptr %1, null
  br i1 %tobool, label %cond.true, label %cond.false

cond.true:                                        ; preds = %entry
  %2 = load ptr, ptr %__begin, align 8
  %call = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE11_M_node_ptrEv(ptr noundef nonnull align 8 dereferenceable(48) %2) #10
  br label %cond.end

cond.false:                                       ; preds = %entry
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi ptr [ %call, %cond.true ], [ null, %cond.false ]
  ret ptr %cond
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt18_Rb_tree_node_base10_S_minimumEPS_(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  br label %while.cond

while.cond:                                       ; preds = %while.body, %entry
  %0 = load ptr, ptr %__x.addr, align 8
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %0, i32 0, i32 2
  %1 = load ptr, ptr %_M_left, align 8
  %cmp = icmp ne ptr %1, null
  br i1 %cmp, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %2 = load ptr, ptr %__x.addr, align 8
  %_M_left1 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %2, i32 0, i32 2
  %3 = load ptr, ptr %_M_left1, align 8
  store ptr %3, ptr %__x.addr, align 8
  br label %while.cond, !llvm.loop !20

while.end:                                        ; preds = %while.cond
  %4 = load ptr, ptr %__x.addr, align 8
  ret ptr %4
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_leftmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 2
  ret ptr %_M_left
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt18_Rb_tree_node_base10_S_maximumEPS_(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  br label %while.cond

while.cond:                                       ; preds = %while.body, %entry
  %0 = load ptr, ptr %__x.addr, align 8
  %_M_right = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %0, i32 0, i32 3
  %1 = load ptr, ptr %_M_right, align 8
  %cmp = icmp ne ptr %1, null
  br i1 %cmp, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %2 = load ptr, ptr %__x.addr, align 8
  %_M_right1 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %2, i32 0, i32 3
  %3 = load ptr, ptr %_M_right1, align 8
  store ptr %3, ptr %__x.addr, align 8
  br label %while.cond, !llvm.loop !21

while.end:                                        ; preds = %while.cond
  %4 = load ptr, ptr %__x.addr, align 8
  ret ptr %4
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_rightmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_right = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 3
  ret ptr %_M_right
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_clone_nodeILb0ENSB_11_Alloc_nodeEEEPSt13_Rb_tree_nodeIS5_ESG_RT0_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__x, ptr noundef nonnull align 8 dereferenceable(8) %__node_gen) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__node_gen.addr = alloca ptr, align 8
  %__tmp = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__node_gen, ptr %__node_gen.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__node_gen.addr, align 8, !nonnull !10, !align !11
  %1 = load ptr, ptr %__x.addr, align 8
  %call = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %1)
  %call2 = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_Alloc_nodeclIRKS5_EEPSt13_Rb_tree_nodeIS5_EOT_(ptr noundef nonnull align 8 dereferenceable(8) %0, ptr noundef nonnull align 8 dereferenceable(12) %call)
  store ptr %call2, ptr %__tmp, align 8
  %2 = load ptr, ptr %__x.addr, align 8
  %_M_color = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %2, i32 0, i32 0
  %3 = load i32, ptr %_M_color, align 8
  %4 = load ptr, ptr %__tmp, align 8
  %_M_color3 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %4, i32 0, i32 0
  store i32 %3, ptr %_M_color3, align 8
  %5 = load ptr, ptr %__tmp, align 8
  %_M_right = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %5, i32 0, i32 3
  store ptr null, ptr %_M_right, align 8
  %6 = load ptr, ptr %__tmp, align 8
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %6, i32 0, i32 2
  store ptr null, ptr %_M_left, align 8
  %7 = load ptr, ptr %__tmp, align 8
  ret ptr %7
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt13_Rb_tree_nodeIS5_E(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %_M_right = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %0, i32 0, i32 3
  %1 = load ptr, ptr %_M_right, align 8
  %tobool = icmp ne ptr %1, null
  br i1 %tobool, label %cond.true, label %cond.false

cond.true:                                        ; preds = %entry
  %2 = load ptr, ptr %__x.addr, align 8
  %_M_right1 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %2, i32 0, i32 3
  %3 = load ptr, ptr %_M_right1, align 8
  %call = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE11_M_node_ptrEv(ptr noundef nonnull align 8 dereferenceable(48) %3) #10
  br label %cond.end

cond.false:                                       ; preds = %entry
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi ptr [ %call, %cond.true ], [ null, %cond.false ]
  ret ptr %cond
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt13_Rb_tree_nodeIS5_E(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %0, i32 0, i32 2
  %1 = load ptr, ptr %_M_left, align 8
  %tobool = icmp ne ptr %1, null
  br i1 %tobool, label %cond.true, label %cond.false

cond.true:                                        ; preds = %entry
  %2 = load ptr, ptr %__x.addr, align 8
  %_M_left1 = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %2, i32 0, i32 2
  %3 = load ptr, ptr %_M_left1, align 8
  %call = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE11_M_node_ptrEv(ptr noundef nonnull align 8 dereferenceable(48) %3) #10
  br label %cond.end

cond.false:                                       ; preds = %entry
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi ptr [ %call, %cond.true ], [ null, %cond.false ]
  ret ptr %cond
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_Alloc_nodeclIRKS5_EEPSt13_Rb_tree_nodeIS5_EOT_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(12) %__arg) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__arg.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__arg, ptr %__arg.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Alloc_node", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %_M_t, align 8, !nonnull !10, !align !11
  %1 = load ptr, ptr %__arg.addr, align 8, !nonnull !10, !align !11
  %call = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_create_nodeIJRKS5_EEEPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %0, ptr noundef nonnull align 8 dereferenceable(12) %1)
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_storage = getelementptr inbounds nuw %"struct.std::_Rb_tree_node", ptr %this1, i32 0, i32 1
  %call = call noundef ptr @_ZN9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE6_M_ptrEv(ptr noundef nonnull align 8 dereferenceable(16) %_M_storage) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_create_nodeIJRKS5_EEEPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(12) %__args) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__args.addr = alloca ptr, align 8
  %__tmp = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__args, ptr %__args.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_get_nodeEv(ptr noundef nonnull align 8 dereferenceable(48) %this1)
  store ptr %call, ptr %__tmp, align 8
  %0 = load ptr, ptr %__tmp, align 8
  %1 = load ptr, ptr %__args.addr, align 8, !nonnull !10, !align !11
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE17_M_construct_nodeIJRKS5_EEEvPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %0, ptr noundef nonnull align 8 dereferenceable(12) %1)
  %2 = load ptr, ptr %__tmp, align 8
  ret ptr %2
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_get_nodeEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %__a.addr.i = alloca ptr, align 8
  %__n.addr.i = alloca i64, align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef nonnull align 1 dereferenceable(1) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE21_M_get_Node_allocatorEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  store ptr %call, ptr %__a.addr.i, align 8
  store i64 1, ptr %__n.addr.i, align 8
  %0 = load ptr, ptr %__a.addr.i, align 8, !nonnull !10
  %1 = load i64, ptr %__n.addr.i, align 8
  %call.i = call noundef ptr @_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEE8allocateEmPKv(ptr noundef nonnull align 1 dereferenceable(1) %0, i64 noundef %1, ptr noundef null)
  ret ptr %call.i
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE17_M_construct_nodeIJRKS5_EEEvPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__node, ptr noundef nonnull align 8 dereferenceable(12) %__args) #1 comdat align 2 {
entry:
  %this.addr.i = alloca ptr, align 8
  %__p.addr.i3 = alloca ptr, align 8
  %__args.addr.i4 = alloca ptr, align 8
  %__a.addr.i = alloca ptr, align 8
  %__p.addr.i = alloca ptr, align 8
  %__args.addr.i = alloca ptr, align 8
  %this.addr = alloca ptr, align 8
  %__node.addr = alloca ptr, align 8
  %__args.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__node, ptr %__node.addr, align 8
  store ptr %__args, ptr %__args.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__node.addr, align 8
  %call = call noundef nonnull align 1 dereferenceable(1) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE21_M_get_Node_allocatorEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %1 = load ptr, ptr %__node.addr, align 8
  %call2 = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %1)
  %2 = load ptr, ptr %__args.addr, align 8, !nonnull !10, !align !11
  store ptr %call, ptr %__a.addr.i, align 8
  store ptr %call2, ptr %__p.addr.i, align 8
  store ptr %2, ptr %__args.addr.i, align 8
  %3 = load ptr, ptr %__a.addr.i, align 8, !nonnull !10
  %4 = load ptr, ptr %__p.addr.i, align 8
  %5 = load ptr, ptr %__args.addr.i, align 8, !nonnull !10, !align !11
  store ptr %3, ptr %this.addr.i, align 8
  store ptr %4, ptr %__p.addr.i3, align 8
  store ptr %5, ptr %__args.addr.i4, align 8
  %this1.i = load ptr, ptr %this.addr.i, align 8
  %6 = load ptr, ptr %__p.addr.i3, align 8
  %7 = load ptr, ptr %__args.addr.i4, align 8, !nonnull !10, !align !11
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %6, ptr align 8 %7, i64 16, i1 false)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 1 dereferenceable(1) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE21_M_get_Node_allocatorEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  ret ptr %_M_impl
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEE8allocateEmPKv(ptr noundef nonnull align 1 dereferenceable(1) %this, i64 noundef %__n, ptr noundef %0) #1 comdat align 2 {
entry:
  %this.addr.i = alloca ptr, align 8
  %this.addr = alloca ptr, align 8
  %__n.addr = alloca i64, align 8
  %.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store i64 %__n, ptr %__n.addr, align 8
  store ptr %0, ptr %.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %1 = load i64, ptr %__n.addr, align 8
  store ptr %this1, ptr %this.addr.i, align 8
  %this1.i = load ptr, ptr %this.addr.i, align 8
  %cmp = icmp ugt i64 %1, 192153584101141162
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  %2 = load i64, ptr %__n.addr, align 8
  %cmp2 = icmp ugt i64 %2, 384307168202282325
  br i1 %cmp2, label %if.then3, label %if.end

if.then3:                                         ; preds = %if.then
  call void @_ZSt28__throw_bad_array_new_lengthv() #12
  unreachable

if.end:                                           ; preds = %if.then
  call void @_ZSt17__throw_bad_allocv() #12
  unreachable

if.else:                                          ; preds = %entry
  %3 = load i64, ptr %__n.addr, align 8
  %mul = mul i64 %3, 48
  %call4 = call noalias noundef nonnull ptr @_Znwm(i64 noundef %mul) #13
  ret ptr %call4
}

; Function Attrs: noreturn
declare void @_ZSt28__throw_bad_array_new_lengthv() #5

; Function Attrs: noreturn
declare void @_ZSt17__throw_bad_allocv() #5

; Function Attrs: nobuiltin allocsize(0)
declare noalias noundef nonnull ptr @_Znwm(i64 noundef) #6

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZN9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE6_M_ptrEv(ptr noundef nonnull align 8 dereferenceable(16) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZN9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE7_M_addrEv(ptr noundef nonnull align 8 dereferenceable(16) %this1) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZN9__gnu_cxx16__aligned_membufISt4pairIKPK1VjEE7_M_addrEv(ptr noundef nonnull align 8 dereferenceable(16) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_storage = getelementptr inbounds nuw %"struct.__gnu_cxx::__aligned_membuf", ptr %this1, i32 0, i32 0
  ret ptr %_M_storage
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE11_M_node_ptrEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  ret ptr %this1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_M_begin_nodeEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_eraseEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %call)
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %_M_impl) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_eraseEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__x) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__y = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  br label %while.cond

while.cond:                                       ; preds = %while.body, %entry
  %0 = load ptr, ptr %__x.addr, align 8
  %tobool = icmp ne ptr %0, null
  br i1 %tobool, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %1 = load ptr, ptr %__x.addr, align 8
  %call = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt13_Rb_tree_nodeIS5_E(ptr noundef %1) #10
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_eraseEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %call)
  %2 = load ptr, ptr %__x.addr, align 8
  %call2 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt13_Rb_tree_nodeIS5_E(ptr noundef %2)
  store ptr %call2, ptr %__y, align 8
  %3 = load ptr, ptr %__x.addr, align 8
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_drop_nodeEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %3) #10
  %4 = load ptr, ptr %__y, align 8
  store ptr %4, ptr %__x.addr, align 8
  br label %while.cond, !llvm.loop !22

while.end:                                        ; preds = %while.cond
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE13_Rb_tree_implIS9_Lb1EED2Ev(ptr noundef nonnull align 8 dereferenceable(48) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr.i = alloca ptr, align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  store ptr %this1, ptr %this.addr.i, align 8
  %this1.i = load ptr, ptr %this.addr.i, align 8
  call void @_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEED2Ev(ptr noundef nonnull align 1 dereferenceable(1) %this1.i) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_drop_nodeEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__p) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__p.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__p, ptr %__p.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__p.addr, align 8
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE15_M_destroy_nodeEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %0) #10
  %1 = load ptr, ptr %__p.addr, align 8
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_put_nodeEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %1) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE15_M_destroy_nodeEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__p) #1 comdat align 2 {
entry:
  %this.addr.i = alloca ptr, align 8
  %__p.addr.i3 = alloca ptr, align 8
  %__a.addr.i = alloca ptr, align 8
  %__p.addr.i = alloca ptr, align 8
  %this.addr = alloca ptr, align 8
  %__p.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__p, ptr %__p.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef nonnull align 1 dereferenceable(1) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE21_M_get_Node_allocatorEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %0 = load ptr, ptr %__p.addr, align 8
  %call2 = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %0)
  store ptr %call, ptr %__a.addr.i, align 8
  store ptr %call2, ptr %__p.addr.i, align 8
  %1 = load ptr, ptr %__a.addr.i, align 8, !nonnull !10
  %2 = load ptr, ptr %__p.addr.i, align 8
  store ptr %1, ptr %this.addr.i, align 8
  store ptr %2, ptr %__p.addr.i3, align 8
  %this1.i = load ptr, ptr %this.addr.i, align 8
  %3 = load ptr, ptr %__p.addr.i3, align 8
  %4 = load ptr, ptr %__p.addr, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_put_nodeEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__p) #1 comdat align 2 {
entry:
  %__a.addr.i = alloca ptr, align 8
  %__p.addr.i = alloca ptr, align 8
  %__n.addr.i = alloca i64, align 8
  %this.addr = alloca ptr, align 8
  %__p.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__p, ptr %__p.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef nonnull align 1 dereferenceable(1) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE21_M_get_Node_allocatorEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %0 = load ptr, ptr %__p.addr, align 8
  store ptr %call, ptr %__a.addr.i, align 8
  store ptr %0, ptr %__p.addr.i, align 8
  store i64 1, ptr %__n.addr.i, align 8
  %1 = load ptr, ptr %__a.addr.i, align 8, !nonnull !10
  %2 = load ptr, ptr %__p.addr.i, align 8
  %3 = load i64, ptr %__n.addr.i, align 8
  call void @_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEE10deallocateEPS7_m(ptr noundef nonnull align 1 dereferenceable(1) %1, ptr noundef %2, i64 noundef %3)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt15__new_allocatorISt13_Rb_tree_nodeISt4pairIKPK1VjEEE10deallocateEPS7_m(ptr noundef nonnull align 1 dereferenceable(1) %this, ptr noundef %__p, i64 noundef %__n) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__p.addr = alloca ptr, align 8
  %__n.addr = alloca i64, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__p, ptr %__p.addr, align 8
  store i64 %__n, ptr %__n.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__p.addr, align 8
  %1 = load i64, ptr %__n.addr, align 8
  %mul = mul i64 %1, 48
  call void @_ZdlPvm(ptr noundef %0, i64 noundef %mul) #14
  ret void
}

; Function Attrs: nobuiltin nounwind
declare void @_ZdlPvm(ptr noundef, i64 noundef) #7

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE11lower_boundERS6_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(8) %__x) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %call = call ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11lower_boundERS4_(ptr noundef nonnull align 8 dereferenceable(48) %_M_t, ptr noundef nonnull align 8 dereferenceable(8) %0)
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %coerce.dive2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %1 = load ptr, ptr %coerce.dive2, align 8
  ret ptr %1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef zeroext i1 @_ZSteqRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y) #1 comdat {
entry:
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %0, i32 0, i32 0
  %1 = load ptr, ptr %_M_node, align 8
  %2 = load ptr, ptr %__y.addr, align 8, !nonnull !10, !align !11
  %_M_node1 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %2, i32 0, i32 0
  %3 = load ptr, ptr %_M_node1, align 8
  %cmp = icmp eq ptr %1, %3
  ret i1 %cmp
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNKSt3mapIPK1VjSt4lessIS2_ESaISt4pairIKS2_jEEE8key_compEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %undef.agg.tmp = alloca %"struct.std::less", align 1
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"class.std::map", ptr %this1, i32 0, i32 0
  call void @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8key_compEv(ptr noundef nonnull align 8 dereferenceable(48) %_M_t)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef zeroext i1 @_ZNKSt4lessIPK1VEclES2_S2_(ptr noundef nonnull align 1 dereferenceable(1) %this, ptr noundef %__x, ptr noundef %__y) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %1 = ptrtoint ptr %0 to i64
  %2 = load ptr, ptr %__y.addr, align 8
  %3 = ptrtoint ptr %2 to i64
  %cmp = icmp ult i64 %1, %3
  ret i1 %cmp
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(12) ptr @_ZNKSt17_Rb_tree_iteratorISt4pairIKPK1VjEEdeEv(ptr noundef nonnull align 8 dereferenceable(8) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %0)
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE22_M_emplace_hint_uniqueIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEESt17_Rb_tree_iteratorIS5_ESt23_Rb_tree_const_iteratorIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr %__pos.coerce, ptr noundef nonnull align 1 dereferenceable(1) %__args, ptr noundef nonnull align 8 dereferenceable(8) %__args1, ptr noundef nonnull align 1 dereferenceable(1) %__args3) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %__pos = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %this.addr = alloca ptr, align 8
  %__args.addr = alloca ptr, align 8
  %__args.addr2 = alloca ptr, align 8
  %__args.addr4 = alloca ptr, align 8
  %__z = alloca %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", align 8
  %__res = alloca %"struct.std::pair.1", align 8
  %agg.tmp = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %agg.tmp8 = alloca %"struct.std::pair.1", align 8
  %cleanup.dest.slot = alloca i32, align 4
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__pos, i32 0, i32 0
  store ptr %__pos.coerce, ptr %coerce.dive, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__args, ptr %__args.addr, align 8
  store ptr %__args1, ptr %__args.addr2, align 8
  store ptr %__args3, ptr %__args.addr4, align 8
  %this5 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__args.addr, align 8, !nonnull !10
  %1 = load ptr, ptr %__args.addr2, align 8, !nonnull !10, !align !11
  %2 = load ptr, ptr %__args.addr4, align 8, !nonnull !10
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_nodeC2IJRKSt21piecewise_construct_tSt5tupleIJRS4_EESH_IJEEEEERSB_DpOT_(ptr noundef nonnull align 8 dereferenceable(16) %__z, ptr noundef nonnull align 8 dereferenceable(48) %this5, ptr noundef nonnull align 1 dereferenceable(1) %0, ptr noundef nonnull align 8 dereferenceable(8) %1, ptr noundef nonnull align 1 dereferenceable(1) %2)
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %agg.tmp, ptr align 8 %__pos, i64 8, i1 false)
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_node6_M_keyEv(ptr noundef nonnull align 8 dereferenceable(16) %__z)
  %coerce.dive6 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %agg.tmp, i32 0, i32 0
  %3 = load ptr, ptr %coerce.dive6, align 8
  %call7 = call { ptr, ptr } @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE29_M_get_insert_hint_unique_posESt23_Rb_tree_const_iteratorIS5_ERS4_(ptr noundef nonnull align 8 dereferenceable(48) %this5, ptr %3, ptr noundef nonnull align 8 dereferenceable(8) %call)
  %4 = getelementptr inbounds nuw { ptr, ptr }, ptr %__res, i32 0, i32 0
  %5 = extractvalue { ptr, ptr } %call7, 0
  store ptr %5, ptr %4, align 8
  %6 = getelementptr inbounds nuw { ptr, ptr }, ptr %__res, i32 0, i32 1
  %7 = extractvalue { ptr, ptr } %call7, 1
  store ptr %7, ptr %6, align 8
  %second = getelementptr inbounds nuw %"struct.std::pair.1", ptr %__res, i32 0, i32 1
  %8 = load ptr, ptr %second, align 8
  %tobool = icmp ne ptr %8, null
  br i1 %tobool, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %agg.tmp8, ptr align 8 %__res, i64 16, i1 false)
  %9 = getelementptr inbounds nuw { ptr, ptr }, ptr %agg.tmp8, i32 0, i32 0
  %10 = load ptr, ptr %9, align 8
  %11 = getelementptr inbounds nuw { ptr, ptr }, ptr %agg.tmp8, i32 0, i32 1
  %12 = load ptr, ptr %11, align 8
  %call9 = call ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_node9_M_insertES3_IPSt18_Rb_tree_node_baseSE_E(ptr noundef nonnull align 8 dereferenceable(16) %__z, ptr %10, ptr %12)
  %coerce.dive10 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  store ptr %call9, ptr %coerce.dive10, align 8
  store i32 1, ptr %cleanup.dest.slot, align 4
  br label %cleanup

if.end:                                           ; preds = %entry
  %first = getelementptr inbounds nuw %"struct.std::pair.1", ptr %__res, i32 0, i32 0
  %13 = load ptr, ptr %first, align 8
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %retval, ptr noundef %13) #10
  store i32 1, ptr %cleanup.dest.slot, align 4
  br label %cleanup

cleanup:                                          ; preds = %if.end, %if.then
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_nodeD2Ev(ptr noundef nonnull align 8 dereferenceable(16) %__z) #10
  %coerce.dive11 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %14 = load ptr, ptr %coerce.dive11, align 8
  ret ptr %14
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt23_Rb_tree_const_iteratorISt4pairIKPK1VjEEC2ERKSt17_Rb_tree_iteratorIS5_E(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(8) %__it) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__it.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__it, ptr %__it.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__it.addr, align 8, !nonnull !10, !align !11
  %_M_node2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %0, i32 0, i32 0
  %1 = load ptr, ptr %_M_node2, align 8
  store ptr %1, ptr %_M_node, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt5tupleIJRKPK1VEEC2ILb1ETnNSt9enable_ifIXclsr4_TCCIXT_EEE29__is_implicitly_constructibleIS4_EEEbE4typeELb1EEES4_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(8) %__elements) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__elements.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__elements, ptr %__elements.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__elements.addr, align 8, !nonnull !10, !align !11
  call void @_ZNSt11_Tuple_implILm0EJRKPK1VEEC2ES4_(ptr noundef nonnull align 8 dereferenceable(8) %this1, ptr noundef nonnull align 8 dereferenceable(8) %0)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11lower_boundERS4_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(8) %__k) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %this.addr = alloca ptr, align 8
  %__k.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__k, ptr %__k.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %call2 = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %0 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call3 = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_lower_boundEPSt18_Rb_tree_node_baseSD_RS4_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef %call, ptr noundef %call2, ptr noundef nonnull align 8 dereferenceable(8) %0)
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %retval, ptr noundef %call3) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %1 = load ptr, ptr %coerce.dive, align 8
  ret ptr %1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_lower_boundEPSt18_Rb_tree_node_baseSD_RS4_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__x, ptr noundef %__y, ptr noundef nonnull align 8 dereferenceable(8) %__k) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  %__k.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  store ptr %__k, ptr %__k.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  br label %while.cond

while.cond:                                       ; preds = %if.end, %entry
  %0 = load ptr, ptr %__x.addr, align 8
  %tobool = icmp ne ptr %0, null
  br i1 %tobool, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %1 = load ptr, ptr %__x.addr, align 8
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %1)
  %2 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call2 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %call, ptr noundef nonnull align 8 dereferenceable(8) %2)
  br i1 %call2, label %if.else, label %if.then

if.then:                                          ; preds = %while.body
  %3 = load ptr, ptr %__x.addr, align 8
  store ptr %3, ptr %__y.addr, align 8
  %4 = load ptr, ptr %__x.addr, align 8
  %call3 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt18_Rb_tree_node_base(ptr noundef %4) #10
  store ptr %call3, ptr %__x.addr, align 8
  br label %if.end

if.else:                                          ; preds = %while.body
  %5 = load ptr, ptr %__x.addr, align 8
  %call4 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt18_Rb_tree_node_base(ptr noundef %5) #10
  store ptr %call4, ptr %__x.addr, align 8
  br label %if.end

if.end:                                           ; preds = %if.else, %if.then
  br label %while.cond, !llvm.loop !23

while.end:                                        ; preds = %while.cond
  %6 = load ptr, ptr %__y.addr, align 8
  ret ptr %6
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_parent = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 1
  %0 = load ptr, ptr %_M_parent, align 8
  ret ptr %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef %__x) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__x.addr, align 8
  store ptr %0, ptr %_M_node, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(8) %__k1, ptr noundef nonnull align 8 dereferenceable(8) %__k2) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__k1.addr = alloca ptr, align 8
  %__k2.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__k1, ptr %__k1.addr, align 8
  store ptr %__k2, ptr %__k2.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__k1.addr, align 8, !nonnull !10, !align !11
  %1 = load ptr, ptr %0, align 8
  %2 = load ptr, ptr %__k2.addr, align 8, !nonnull !10, !align !11
  %3 = load ptr, ptr %2, align 8
  %call = call noundef zeroext i1 @_ZNKSt4lessIPK1VEclES2_S2_(ptr noundef nonnull align 1 dereferenceable(1) %_M_impl, ptr noundef %1, ptr noundef %3) #10
  ret i1 %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyERKSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %0)
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt18_Rb_tree_node_base(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %0, i32 0, i32 2
  %1 = load ptr, ptr %_M_left, align 8
  ret ptr %1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt18_Rb_tree_node_base(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %_M_right = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %0, i32 0, i32 3
  %1 = load ptr, ptr %_M_right, align 8
  ret ptr %1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyERKSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %__node) #1 comdat align 2 {
entry:
  %__node.addr = alloca ptr, align 8
  %ref.tmp = alloca %"struct.std::_Select1st", align 1
  store ptr %__node, ptr %__node.addr, align 8
  %0 = load ptr, ptr %__node.addr, align 8, !nonnull !10, !align !11
  %call = call noundef ptr @_ZNKSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %0)
  %call1 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNKSt10_Select1stISt4pairIKPK1VjEEclERKS5_(ptr noundef nonnull align 1 dereferenceable(1) %ref.tmp, ptr noundef nonnull align 8 dereferenceable(12) %call)
  ret ptr %call1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNKSt10_Select1stISt4pairIKPK1VjEEclERKS5_(ptr noundef nonnull align 1 dereferenceable(1) %this, ptr noundef nonnull align 8 dereferenceable(12) %__x) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %first = getelementptr inbounds nuw %"struct.std::pair", ptr %0, i32 0, i32 0
  ret ptr %first
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8key_compEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_nodeC2IJRKSt21piecewise_construct_tSt5tupleIJRS4_EESH_IJEEEEERSB_DpOT_(ptr noundef nonnull align 8 dereferenceable(16) %this, ptr noundef nonnull align 8 dereferenceable(48) %__t, ptr noundef nonnull align 1 dereferenceable(1) %__args, ptr noundef nonnull align 8 dereferenceable(8) %__args1, ptr noundef nonnull align 1 dereferenceable(1) %__args3) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__t.addr = alloca ptr, align 8
  %__args.addr = alloca ptr, align 8
  %__args.addr2 = alloca ptr, align 8
  %__args.addr4 = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__t, ptr %__t.addr, align 8
  store ptr %__args, ptr %__args.addr, align 8
  store ptr %__args1, ptr %__args.addr2, align 8
  store ptr %__args3, ptr %__args.addr4, align 8
  %this5 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this5, i32 0, i32 0
  %0 = load ptr, ptr %__t.addr, align 8, !nonnull !10, !align !11
  store ptr %0, ptr %_M_t, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this5, i32 0, i32 1
  %1 = load ptr, ptr %__t.addr, align 8, !nonnull !10, !align !11
  %2 = load ptr, ptr %__args.addr, align 8, !nonnull !10
  %3 = load ptr, ptr %__args.addr2, align 8, !nonnull !10, !align !11
  %4 = load ptr, ptr %__args.addr4, align 8, !nonnull !10
  %call = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_create_nodeIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEEPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %1, ptr noundef nonnull align 1 dereferenceable(1) %2, ptr noundef nonnull align 8 dereferenceable(8) %3, ptr noundef nonnull align 1 dereferenceable(1) %4)
  store ptr %call, ptr %_M_node, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local { ptr, ptr } @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE29_M_get_insert_hint_unique_posESt23_Rb_tree_const_iteratorIS5_ERS4_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr %__position.coerce, ptr noundef nonnull align 8 dereferenceable(8) %__k) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::pair.1", align 8
  %__position = alloca %"struct.std::_Rb_tree_const_iterator", align 8
  %this.addr = alloca ptr, align 8
  %__k.addr = alloca ptr, align 8
  %ref.tmp = alloca ptr, align 8
  %__before = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp31 = alloca ptr, align 8
  %__after = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp48 = alloca ptr, align 8
  %ref.tmp60 = alloca ptr, align 8
  %ref.tmp69 = alloca ptr, align 8
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  store ptr %__position.coerce, ptr %coerce.dive, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__k, ptr %__k.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %cmp = icmp eq ptr %0, %call
  br i1 %cmp, label %if.then, label %if.else10

if.then:                                          ; preds = %entry
  %call2 = call noundef i64 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE4sizeEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %cmp3 = icmp ugt i64 %call2, 0
  br i1 %cmp3, label %land.lhs.true, label %if.else

land.lhs.true:                                    ; preds = %if.then
  %call4 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_rightmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %1 = load ptr, ptr %call4, align 8
  %call5 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %1)
  %2 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call6 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %call5, ptr noundef nonnull align 8 dereferenceable(8) %2)
  br i1 %call6, label %if.then7, label %if.else

if.then7:                                         ; preds = %land.lhs.true
  store ptr null, ptr %ref.tmp, align 8
  %call8 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_rightmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IS1_RS1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp, ptr noundef nonnull align 8 dereferenceable(8) %call8)
  br label %return

if.else:                                          ; preds = %land.lhs.true, %if.then
  %3 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call9 = call { ptr, ptr } @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE24_M_get_insert_unique_posERS4_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %3)
  %4 = getelementptr inbounds nuw { ptr, ptr }, ptr %retval, i32 0, i32 0
  %5 = extractvalue { ptr, ptr } %call9, 0
  store ptr %5, ptr %4, align 8
  %6 = getelementptr inbounds nuw { ptr, ptr }, ptr %retval, i32 0, i32 1
  %7 = extractvalue { ptr, ptr } %call9, 1
  store ptr %7, ptr %6, align 8
  br label %return

if.else10:                                        ; preds = %entry
  %8 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %_M_node11 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %9 = load ptr, ptr %_M_node11, align 8
  %call12 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %9)
  %call13 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %8, ptr noundef nonnull align 8 dereferenceable(8) %call12)
  br i1 %call13, label %if.then14, label %if.else38

if.then14:                                        ; preds = %if.else10
  %_M_node15 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %10 = load ptr, ptr %_M_node15, align 8
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %__before, ptr noundef %10) #10
  %_M_node16 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %11 = load ptr, ptr %_M_node16, align 8
  %call17 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_leftmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %12 = load ptr, ptr %call17, align 8
  %cmp18 = icmp eq ptr %11, %12
  br i1 %cmp18, label %if.then19, label %if.else22

if.then19:                                        ; preds = %if.then14
  %call20 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_leftmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %call21 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_leftmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S4_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %call20, ptr noundef nonnull align 8 dereferenceable(8) %call21)
  br label %return

if.else22:                                        ; preds = %if.then14
  %call23 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEmmEv(ptr noundef nonnull align 8 dereferenceable(8) %__before) #10
  %_M_node24 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %call23, i32 0, i32 0
  %13 = load ptr, ptr %_M_node24, align 8
  %call25 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %13)
  %14 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call26 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %call25, ptr noundef nonnull align 8 dereferenceable(8) %14)
  br i1 %call26, label %if.then27, label %if.else36

if.then27:                                        ; preds = %if.else22
  %_M_node28 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %__before, i32 0, i32 0
  %15 = load ptr, ptr %_M_node28, align 8
  %call29 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt18_Rb_tree_node_base(ptr noundef %15) #10
  %tobool = icmp ne ptr %call29, null
  br i1 %tobool, label %if.else33, label %if.then30

if.then30:                                        ; preds = %if.then27
  store ptr null, ptr %ref.tmp31, align 8
  %_M_node32 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %__before, i32 0, i32 0
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IS1_RS1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp31, ptr noundef nonnull align 8 dereferenceable(8) %_M_node32)
  br label %return

if.else33:                                        ; preds = %if.then27
  %_M_node34 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %_M_node35 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S4_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %_M_node34, ptr noundef nonnull align 8 dereferenceable(8) %_M_node35)
  br label %return

if.else36:                                        ; preds = %if.else22
  %16 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call37 = call { ptr, ptr } @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE24_M_get_insert_unique_posERS4_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %16)
  %17 = getelementptr inbounds nuw { ptr, ptr }, ptr %retval, i32 0, i32 0
  %18 = extractvalue { ptr, ptr } %call37, 0
  store ptr %18, ptr %17, align 8
  %19 = getelementptr inbounds nuw { ptr, ptr }, ptr %retval, i32 0, i32 1
  %20 = extractvalue { ptr, ptr } %call37, 1
  store ptr %20, ptr %19, align 8
  br label %return

if.else38:                                        ; preds = %if.else10
  %_M_node39 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %21 = load ptr, ptr %_M_node39, align 8
  %call40 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %21)
  %22 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call41 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %call40, ptr noundef nonnull align 8 dereferenceable(8) %22)
  br i1 %call41, label %if.then42, label %if.else67

if.then42:                                        ; preds = %if.else38
  %_M_node43 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %23 = load ptr, ptr %_M_node43, align 8
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %__after, ptr noundef %23) #10
  %_M_node44 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %24 = load ptr, ptr %_M_node44, align 8
  %call45 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_rightmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %25 = load ptr, ptr %call45, align 8
  %cmp46 = icmp eq ptr %24, %25
  br i1 %cmp46, label %if.then47, label %if.else50

if.then47:                                        ; preds = %if.then42
  store ptr null, ptr %ref.tmp48, align 8
  %call49 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_rightmostEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IS1_RS1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp48, ptr noundef nonnull align 8 dereferenceable(8) %call49)
  br label %return

if.else50:                                        ; preds = %if.then42
  %26 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call51 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEppEv(ptr noundef nonnull align 8 dereferenceable(8) %__after) #10
  %_M_node52 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %call51, i32 0, i32 0
  %27 = load ptr, ptr %_M_node52, align 8
  %call53 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %27)
  %call54 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %26, ptr noundef nonnull align 8 dereferenceable(8) %call53)
  br i1 %call54, label %if.then55, label %if.else65

if.then55:                                        ; preds = %if.else50
  %_M_node56 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  %28 = load ptr, ptr %_M_node56, align 8
  %call57 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt18_Rb_tree_node_base(ptr noundef %28) #10
  %tobool58 = icmp ne ptr %call57, null
  br i1 %tobool58, label %if.else62, label %if.then59

if.then59:                                        ; preds = %if.then55
  store ptr null, ptr %ref.tmp60, align 8
  %_M_node61 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IS1_RS1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp60, ptr noundef nonnull align 8 dereferenceable(8) %_M_node61)
  br label %return

if.else62:                                        ; preds = %if.then55
  %_M_node63 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %__after, i32 0, i32 0
  %_M_node64 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %__after, i32 0, i32 0
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S4_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %_M_node63, ptr noundef nonnull align 8 dereferenceable(8) %_M_node64)
  br label %return

if.else65:                                        ; preds = %if.else50
  %29 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call66 = call { ptr, ptr } @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE24_M_get_insert_unique_posERS4_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %29)
  %30 = getelementptr inbounds nuw { ptr, ptr }, ptr %retval, i32 0, i32 0
  %31 = extractvalue { ptr, ptr } %call66, 0
  store ptr %31, ptr %30, align 8
  %32 = getelementptr inbounds nuw { ptr, ptr }, ptr %retval, i32 0, i32 1
  %33 = extractvalue { ptr, ptr } %call66, 1
  store ptr %33, ptr %32, align 8
  br label %return

if.else67:                                        ; preds = %if.else38
  %_M_node68 = getelementptr inbounds nuw %"struct.std::_Rb_tree_const_iterator", ptr %__position, i32 0, i32 0
  store ptr null, ptr %ref.tmp69, align 8
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %_M_node68, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp69)
  br label %return

return:                                           ; preds = %if.else67, %if.else65, %if.else62, %if.then59, %if.then47, %if.else36, %if.else33, %if.then30, %if.then19, %if.else, %if.then7
  %34 = load { ptr, ptr }, ptr %retval, align 8
  ret { ptr, ptr } %34
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_node6_M_keyEv(ptr noundef nonnull align 8 dereferenceable(16) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this1, i32 0, i32 1
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt13_Rb_tree_nodeIS5_E(ptr noundef %0)
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_node9_M_insertES3_IPSt18_Rb_tree_node_baseSE_E(ptr noundef nonnull align 8 dereferenceable(16) %this, ptr %__p.coerce0, ptr %__p.coerce1) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %__p = alloca %"struct.std::pair.1", align 8
  %this.addr = alloca ptr, align 8
  %0 = getelementptr inbounds nuw { ptr, ptr }, ptr %__p, i32 0, i32 0
  store ptr %__p.coerce0, ptr %0, align 8
  %1 = getelementptr inbounds nuw { ptr, ptr }, ptr %__p, i32 0, i32 1
  store ptr %__p.coerce1, ptr %1, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_t = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this1, i32 0, i32 0
  %2 = load ptr, ptr %_M_t, align 8, !nonnull !10, !align !11
  %first = getelementptr inbounds nuw %"struct.std::pair.1", ptr %__p, i32 0, i32 0
  %3 = load ptr, ptr %first, align 8
  %second = getelementptr inbounds nuw %"struct.std::pair.1", ptr %__p, i32 0, i32 1
  %4 = load ptr, ptr %second, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this1, i32 0, i32 1
  %5 = load ptr, ptr %_M_node, align 8
  %call = call ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_insert_nodeEPSt18_Rb_tree_node_baseSD_PSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %2, ptr noundef %3, ptr noundef %4, ptr noundef %5)
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  store ptr %call, ptr %coerce.dive, align 8
  %_M_node2 = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this1, i32 0, i32 1
  store ptr null, ptr %_M_node2, align 8
  %coerce.dive3 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %6 = load ptr, ptr %coerce.dive3, align 8
  ret ptr %6
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE10_Auto_nodeD2Ev(ptr noundef nonnull align 8 dereferenceable(16) %this) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this1, i32 0, i32 1
  %0 = load ptr, ptr %_M_node, align 8
  %tobool = icmp ne ptr %0, null
  br i1 %tobool, label %if.then, label %if.end

if.then:                                          ; preds = %entry
  %_M_t = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this1, i32 0, i32 0
  %1 = load ptr, ptr %_M_t, align 8, !nonnull !10, !align !11
  %_M_node2 = getelementptr inbounds nuw %"struct.std::_Rb_tree<const V *, std::pair<const V *const, unsigned int>, std::_Select1st<std::pair<const V *const, unsigned int>>, std::less<const V *>>::_Auto_node", ptr %this1, i32 0, i32 1
  %2 = load ptr, ptr %_M_node2, align 8
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE12_M_drop_nodeEPSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %1, ptr noundef %2) #10
  br label %if.end

if.end:                                           ; preds = %if.then, %entry
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_create_nodeIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEEPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 1 dereferenceable(1) %__args, ptr noundef nonnull align 8 dereferenceable(8) %__args1, ptr noundef nonnull align 1 dereferenceable(1) %__args3) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__args.addr = alloca ptr, align 8
  %__args.addr2 = alloca ptr, align 8
  %__args.addr4 = alloca ptr, align 8
  %__tmp = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__args, ptr %__args.addr, align 8
  store ptr %__args1, ptr %__args.addr2, align 8
  store ptr %__args3, ptr %__args.addr4, align 8
  %this5 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE11_M_get_nodeEv(ptr noundef nonnull align 8 dereferenceable(48) %this5)
  store ptr %call, ptr %__tmp, align 8
  %0 = load ptr, ptr %__tmp, align 8
  %1 = load ptr, ptr %__args.addr, align 8, !nonnull !10
  %2 = load ptr, ptr %__args.addr2, align 8, !nonnull !10, !align !11
  %3 = load ptr, ptr %__args.addr4, align 8, !nonnull !10
  call void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE17_M_construct_nodeIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEEvPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %this5, ptr noundef %0, ptr noundef nonnull align 1 dereferenceable(1) %1, ptr noundef nonnull align 8 dereferenceable(8) %2, ptr noundef nonnull align 1 dereferenceable(1) %3)
  %4 = load ptr, ptr %__tmp, align 8
  ret ptr %4
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE17_M_construct_nodeIJRKSt21piecewise_construct_tSt5tupleIJRS4_EESG_IJEEEEEvPSt13_Rb_tree_nodeIS5_EDpOT_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__node, ptr noundef nonnull align 1 dereferenceable(1) %__args, ptr noundef nonnull align 8 dereferenceable(8) %__args1, ptr noundef nonnull align 1 dereferenceable(1) %__args3) #1 comdat align 2 {
entry:
  %this.addr.i = alloca ptr, align 8
  %__p.addr.i7 = alloca ptr, align 8
  %__args.addr.i8 = alloca ptr, align 8
  %__args.addr2.i9 = alloca ptr, align 8
  %__args.addr4.i10 = alloca ptr, align 8
  %agg.tmp6.i = alloca %"class.std::tuple", align 8
  %__a.addr.i = alloca ptr, align 8
  %__p.addr.i = alloca ptr, align 8
  %__args.addr.i = alloca ptr, align 8
  %__args.addr2.i = alloca ptr, align 8
  %__args.addr4.i = alloca ptr, align 8
  %this.addr = alloca ptr, align 8
  %__node.addr = alloca ptr, align 8
  %__args.addr = alloca ptr, align 8
  %__args.addr2 = alloca ptr, align 8
  %__args.addr4 = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__node, ptr %__node.addr, align 8
  store ptr %__args, ptr %__args.addr, align 8
  store ptr %__args1, ptr %__args.addr2, align 8
  store ptr %__args3, ptr %__args.addr4, align 8
  %this5 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__node.addr, align 8
  %call = call noundef nonnull align 1 dereferenceable(1) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE21_M_get_Node_allocatorEv(ptr noundef nonnull align 8 dereferenceable(48) %this5) #10
  %1 = load ptr, ptr %__node.addr, align 8
  %call6 = call noundef ptr @_ZNSt13_Rb_tree_nodeISt4pairIKPK1VjEE9_M_valptrEv(ptr noundef nonnull align 8 dereferenceable(48) %1)
  %2 = load ptr, ptr %__args.addr, align 8, !nonnull !10
  %3 = load ptr, ptr %__args.addr2, align 8, !nonnull !10, !align !11
  %4 = load ptr, ptr %__args.addr4, align 8, !nonnull !10
  store ptr %call, ptr %__a.addr.i, align 8
  store ptr %call6, ptr %__p.addr.i, align 8
  store ptr %2, ptr %__args.addr.i, align 8
  store ptr %3, ptr %__args.addr2.i, align 8
  store ptr %4, ptr %__args.addr4.i, align 8
  %5 = load ptr, ptr %__a.addr.i, align 8, !nonnull !10
  %6 = load ptr, ptr %__p.addr.i, align 8
  %7 = load ptr, ptr %__args.addr.i, align 8, !nonnull !10
  %8 = load ptr, ptr %__args.addr2.i, align 8, !nonnull !10, !align !11
  %9 = load ptr, ptr %__args.addr4.i, align 8, !nonnull !10
  store ptr %5, ptr %this.addr.i, align 8
  store ptr %6, ptr %__p.addr.i7, align 8
  store ptr %7, ptr %__args.addr.i8, align 8
  store ptr %8, ptr %__args.addr2.i9, align 8
  store ptr %9, ptr %__args.addr4.i10, align 8
  %this5.i = load ptr, ptr %this.addr.i, align 8
  %10 = load ptr, ptr %__p.addr.i7, align 8
  %11 = load ptr, ptr %__args.addr.i8, align 8, !nonnull !10
  %12 = load ptr, ptr %__args.addr2.i9, align 8, !nonnull !10, !align !11
  call void @_ZNSt5tupleIJRKPK1VEEC2EOS5_(ptr noundef nonnull align 8 dereferenceable(8) %agg.tmp6.i, ptr noundef nonnull align 8 dereferenceable(8) %12) #10
  %13 = load ptr, ptr %__args.addr4.i10, align 8, !nonnull !10
  call void @_ZNSt4pairIKPK1VjEC2IJRS3_EJEEESt21piecewise_construct_tSt5tupleIJDpT_EES8_IJDpT0_EE(ptr noundef nonnull align 8 dereferenceable(12) %10, ptr dead_on_return noundef %agg.tmp6.i)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt5tupleIJRKPK1VEEC2EOS5_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(8) %0) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %0, ptr %.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %1 = load ptr, ptr %.addr, align 8, !nonnull !10, !align !11
  call void @_ZNSt11_Tuple_implILm0EJRKPK1VEEC2EOS5_(ptr noundef nonnull align 8 dereferenceable(8) %this1, ptr noundef nonnull align 8 dereferenceable(8) %1) #10
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt4pairIKPK1VjEC2IJRS3_EJEEESt21piecewise_construct_tSt5tupleIJDpT_EES8_IJDpT0_EE(ptr noundef nonnull align 8 dereferenceable(12) %this, ptr dead_on_return noundef %__first) unnamed_addr #1 comdat align 2 {
entry:
  %0 = alloca %"struct.std::piecewise_construct_t", align 1
  %__second = alloca %"class.std::tuple.0", align 1
  %this.addr = alloca ptr, align 8
  %__first.indirect_addr = alloca ptr, align 8
  %agg.tmp = alloca %"struct.std::_Index_tuple", align 1
  %agg.tmp3 = alloca %"struct.std::_Index_tuple.3", align 1
  store ptr %this, ptr %this.addr, align 8
  store ptr %__first, ptr %__first.indirect_addr, align 8
  %this2 = load ptr, ptr %this.addr, align 8
  call void @_ZNSt4pairIKPK1VjEC2IJRS3_EJLm0EEJETpTnmJEEERSt5tupleIJDpT_EERS7_IJDpT1_EESt12_Index_tupleIJXspT0_EEESG_IJXspT2_EEE(ptr noundef nonnull align 8 dereferenceable(12) %this2, ptr noundef nonnull align 8 dereferenceable(8) %__first, ptr noundef nonnull align 1 dereferenceable(1) %__second)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt11_Tuple_implILm0EJRKPK1VEEC2EOS5_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(8) %__in) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__in.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__in, ptr %__in.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__in.addr, align 8, !nonnull !10, !align !11
  call void @llvm.memcpy.p0.p0.i64(ptr align 8 %this1, ptr align 8 %0, i64 8, i1 false)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt4pairIKPK1VjEC2IJRS3_EJLm0EEJETpTnmJEEERSt5tupleIJDpT_EERS7_IJDpT1_EESt12_Index_tupleIJXspT0_EEESG_IJXspT2_EEE(ptr noundef nonnull align 8 dereferenceable(12) %this, ptr noundef nonnull align 8 dereferenceable(8) %__tuple1, ptr noundef nonnull align 1 dereferenceable(1) %__tuple2) unnamed_addr #1 comdat align 2 {
entry:
  %0 = alloca %"struct.std::_Index_tuple", align 1
  %1 = alloca %"struct.std::_Index_tuple.3", align 1
  %this.addr = alloca ptr, align 8
  %__tuple1.addr = alloca ptr, align 8
  %__tuple2.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__tuple1, ptr %__tuple1.addr, align 8
  store ptr %__tuple2, ptr %__tuple2.addr, align 8
  %this2 = load ptr, ptr %this.addr, align 8
  %first = getelementptr inbounds nuw %"struct.std::pair", ptr %this2, i32 0, i32 0
  %2 = load ptr, ptr %__tuple1.addr, align 8, !nonnull !10, !align !11
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZSt3getILm0EJRKPK1VEERNSt13tuple_elementIXT_ESt5tupleIJDpT0_EEE4typeERS9_(ptr noundef nonnull align 8 dereferenceable(8) %2) #10
  %3 = load ptr, ptr %call, align 8
  store ptr %3, ptr %first, align 8
  %second = getelementptr inbounds nuw %"struct.std::pair", ptr %this2, i32 0, i32 1
  store i32 0, ptr %second, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZSt3getILm0EJRKPK1VEERNSt13tuple_elementIXT_ESt5tupleIJDpT0_EEE4typeERS9_(ptr noundef nonnull align 8 dereferenceable(8) %__t) #1 comdat {
entry:
  %__t.addr = alloca ptr, align 8
  store ptr %__t, ptr %__t.addr, align 8
  %0 = load ptr, ptr %__t.addr, align 8, !nonnull !10, !align !11
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZSt12__get_helperILm0ERKPK1VJEERT0_RSt11_Tuple_implIXT_EJS5_DpT1_EE(ptr noundef nonnull align 8 dereferenceable(8) %0) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZSt12__get_helperILm0ERKPK1VJEERT0_RSt11_Tuple_implIXT_EJS5_DpT1_EE(ptr noundef nonnull align 8 dereferenceable(8) %__t) #1 comdat {
entry:
  %__t.addr = alloca ptr, align 8
  store ptr %__t, ptr %__t.addr, align 8
  %0 = load ptr, ptr %__t.addr, align 8, !nonnull !10, !align !11
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt11_Tuple_implILm0EJRKPK1VEE7_M_headERS5_(ptr noundef nonnull align 8 dereferenceable(8) %0) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt11_Tuple_implILm0EJRKPK1VEE7_M_headERS5_(ptr noundef nonnull align 8 dereferenceable(8) %__t) #1 comdat align 2 {
entry:
  %__t.addr = alloca ptr, align 8
  store ptr %__t, ptr %__t.addr, align 8
  %0 = load ptr, ptr %__t.addr, align 8, !nonnull !10, !align !11
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt10_Head_baseILm0ERKPK1VLb0EE7_M_headERS5_(ptr noundef nonnull align 8 dereferenceable(8) %0) #10
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt10_Head_baseILm0ERKPK1VLb0EE7_M_headERS5_(ptr noundef nonnull align 8 dereferenceable(8) %__b) #1 comdat align 2 {
entry:
  %__b.addr = alloca ptr, align 8
  store ptr %__b, ptr %__b.addr, align 8
  %0 = load ptr, ptr %__b.addr, align 8, !nonnull !10, !align !11
  %_M_head_impl = getelementptr inbounds nuw %"struct.std::_Head_base", ptr %0, i32 0, i32 0
  %1 = load ptr, ptr %_M_head_impl, align 8, !nonnull !10, !align !11
  ret ptr %1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef i64 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE4sizeEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_node_count = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 1
  %0 = load i64, ptr %_M_node_count, align 8
  ret i64 %0
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IS1_RS1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %this, ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %first = getelementptr inbounds nuw %"struct.std::pair.1", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %1 = load ptr, ptr %0, align 8
  store ptr %1, ptr %first, align 8
  %second = getelementptr inbounds nuw %"struct.std::pair.1", ptr %this1, i32 0, i32 1
  %2 = load ptr, ptr %__y.addr, align 8, !nonnull !10, !align !11
  %3 = load ptr, ptr %2, align 8
  store ptr %3, ptr %second, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local { ptr, ptr } @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE24_M_get_insert_unique_posERS4_(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef nonnull align 8 dereferenceable(8) %__k) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::pair.1", align 8
  %this.addr = alloca ptr, align 8
  %__k.addr = alloca ptr, align 8
  %__x = alloca ptr, align 8
  %__y = alloca ptr, align 8
  %__comp = alloca i8, align 1
  %__j = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp = alloca %"struct.std::_Rb_tree_iterator", align 8
  %ref.tmp18 = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__k, ptr %__k.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_M_beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  store ptr %call, ptr %__x, align 8
  %call2 = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  store ptr %call2, ptr %__y, align 8
  store i8 1, ptr %__comp, align 1
  br label %while.cond

while.cond:                                       ; preds = %cond.end, %entry
  %0 = load ptr, ptr %__x, align 8
  %tobool = icmp ne ptr %0, null
  br i1 %tobool, label %while.body, label %while.end

while.body:                                       ; preds = %while.cond
  %1 = load ptr, ptr %__x, align 8
  store ptr %1, ptr %__y, align 8
  %2 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %3 = load ptr, ptr %__x, align 8
  %call3 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %3)
  %call4 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %2, ptr noundef nonnull align 8 dereferenceable(8) %call3)
  %storedv = zext i1 %call4 to i8
  store i8 %storedv, ptr %__comp, align 1
  %4 = load i8, ptr %__comp, align 1
  %loadedv = trunc i8 %4 to i1
  br i1 %loadedv, label %cond.true, label %cond.false

cond.true:                                        ; preds = %while.body
  %5 = load ptr, ptr %__x, align 8
  %call5 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE7_S_leftEPSt18_Rb_tree_node_base(ptr noundef %5) #10
  br label %cond.end

cond.false:                                       ; preds = %while.body
  %6 = load ptr, ptr %__x, align 8
  %call6 = call noundef ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE8_S_rightEPSt18_Rb_tree_node_base(ptr noundef %6) #10
  br label %cond.end

cond.end:                                         ; preds = %cond.false, %cond.true
  %cond = phi ptr [ %call5, %cond.true ], [ %call6, %cond.false ]
  store ptr %cond, ptr %__x, align 8
  br label %while.cond, !llvm.loop !24

while.end:                                        ; preds = %while.cond
  %7 = load ptr, ptr %__y, align 8
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %__j, ptr noundef %7) #10
  %8 = load i8, ptr %__comp, align 1
  %loadedv7 = trunc i8 %8 to i1
  br i1 %loadedv7, label %if.then, label %if.end12

if.then:                                          ; preds = %while.end
  %call8 = call ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %ref.tmp, i32 0, i32 0
  store ptr %call8, ptr %coerce.dive, align 8
  %call9 = call noundef zeroext i1 @_ZSteqRKSt17_Rb_tree_iteratorISt4pairIKPK1VjEES8_(ptr noundef nonnull align 8 dereferenceable(8) %__j, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp) #10
  br i1 %call9, label %if.then10, label %if.else

if.then10:                                        ; preds = %if.then
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S4_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y)
  br label %return

if.else:                                          ; preds = %if.then
  %call11 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEmmEv(ptr noundef nonnull align 8 dereferenceable(8) %__j) #10
  br label %if.end

if.end:                                           ; preds = %if.else
  br label %if.end12

if.end12:                                         ; preds = %if.end, %while.end
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %__j, i32 0, i32 0
  %9 = load ptr, ptr %_M_node, align 8
  %call13 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %9)
  %10 = load ptr, ptr %__k.addr, align 8, !nonnull !10, !align !11
  %call14 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %call13, ptr noundef nonnull align 8 dereferenceable(8) %10)
  br i1 %call14, label %if.then15, label %if.end16

if.then15:                                        ; preds = %if.end12
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S4_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y)
  br label %return

if.end16:                                         ; preds = %if.end12
  %_M_node17 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %__j, i32 0, i32 0
  store ptr null, ptr %ref.tmp18, align 8
  call void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %retval, ptr noundef nonnull align 8 dereferenceable(8) %_M_node17, ptr noundef nonnull align 8 dereferenceable(8) %ref.tmp18)
  br label %return

return:                                           ; preds = %if.end16, %if.then15, %if.then10
  %11 = load { ptr, ptr }, ptr %retval, align 8
  ret { ptr, ptr } %11
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S4_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %this, ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %first = getelementptr inbounds nuw %"struct.std::pair.1", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %1 = load ptr, ptr %0, align 8
  store ptr %1, ptr %first, align 8
  %second = getelementptr inbounds nuw %"struct.std::pair.1", ptr %this1, i32 0, i32 1
  %2 = load ptr, ptr %__y.addr, align 8, !nonnull !10, !align !11
  %3 = load ptr, ptr %2, align 8
  store ptr %3, ptr %second, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEmmEv(ptr noundef nonnull align 8 dereferenceable(8) %this) #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_node = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %_M_node, align 8
  %call = call noundef ptr @_ZSt18_Rb_tree_decrementPSt18_Rb_tree_node_base(ptr noundef %0) #11
  %_M_node2 = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %this1, i32 0, i32 0
  store ptr %call, ptr %_M_node2, align 8
  ret ptr %this1
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt4pairIPSt18_Rb_tree_node_baseS1_EC2IRS1_S1_TnNSt9enable_ifIXaaclsr5_PCCPE22_MoveConstructiblePairIT_T0_EEclsr5_PCCPE30_ImplicitlyMoveConvertiblePairIS6_S7_EEEbE4typeELb1EEEOS6_OS7_(ptr noundef nonnull align 8 dereferenceable(16) %this, ptr noundef nonnull align 8 dereferenceable(8) %__x, ptr noundef nonnull align 8 dereferenceable(8) %__y) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__y.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__y, ptr %__y.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %first = getelementptr inbounds nuw %"struct.std::pair.1", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__x.addr, align 8, !nonnull !10, !align !11
  %1 = load ptr, ptr %0, align 8
  store ptr %1, ptr %first, align 8
  %second = getelementptr inbounds nuw %"struct.std::pair.1", ptr %this1, i32 0, i32 1
  %2 = load ptr, ptr %__y.addr, align 8, !nonnull !10, !align !11
  %3 = load ptr, ptr %2, align 8
  store ptr %3, ptr %second, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE5beginEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  %_M_left = getelementptr inbounds nuw %"struct.std::_Rb_tree_node_base", ptr %_M_header, i32 0, i32 2
  %0 = load ptr, ptr %_M_left, align 8
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %retval, ptr noundef %0) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %1 = load ptr, ptr %coerce.dive, align 8
  ret ptr %1
}

; Function Attrs: nounwind willreturn memory(read)
declare noundef ptr @_ZSt18_Rb_tree_decrementPSt18_Rb_tree_node_base(ptr noundef) #4

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt13_Rb_tree_nodeIS5_E(ptr noundef %__x) #1 comdat align 2 {
entry:
  %__x.addr = alloca ptr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %call = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyERKSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %0)
  ret ptr %call
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_insert_nodeEPSt18_Rb_tree_node_baseSD_PSt13_Rb_tree_nodeIS5_E(ptr noundef nonnull align 8 dereferenceable(48) %this, ptr noundef %__x, ptr noundef %__p, ptr noundef %__z) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %this.addr = alloca ptr, align 8
  %__x.addr = alloca ptr, align 8
  %__p.addr = alloca ptr, align 8
  %__z.addr = alloca ptr, align 8
  %__insert_left = alloca i8, align 1
  %__base_z = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__p, ptr %__p.addr, align 8
  store ptr %__z, ptr %__z.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__x.addr, align 8
  %tobool = icmp ne ptr %0, null
  br i1 %tobool, label %lor.end, label %lor.lhs.false

lor.lhs.false:                                    ; preds = %entry
  %1 = load ptr, ptr %__p.addr, align 8
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  %cmp = icmp eq ptr %1, %call
  br i1 %cmp, label %lor.end, label %lor.rhs

lor.rhs:                                          ; preds = %lor.lhs.false
  %2 = load ptr, ptr %__z.addr, align 8
  %call2 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt13_Rb_tree_nodeIS5_E(ptr noundef %2)
  %3 = load ptr, ptr %__p.addr, align 8
  %call3 = call noundef nonnull align 8 dereferenceable(8) ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_S_keyEPSt18_Rb_tree_node_base(ptr noundef %3)
  %call4 = call noundef zeroext i1 @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE14_M_key_compareIS2_S2_EEbRKT_RKT0_(ptr noundef nonnull align 8 dereferenceable(48) %this1, ptr noundef nonnull align 8 dereferenceable(8) %call2, ptr noundef nonnull align 8 dereferenceable(8) %call3)
  br label %lor.end

lor.end:                                          ; preds = %lor.rhs, %lor.lhs.false, %entry
  %4 = phi i1 [ true, %lor.lhs.false ], [ true, %entry ], [ %call4, %lor.rhs ]
  %storedv = zext i1 %4 to i8
  store i8 %storedv, ptr %__insert_left, align 1
  %5 = load ptr, ptr %__z.addr, align 8
  %call5 = call noundef ptr @_ZNKSt18_Rb_tree_node_base11_M_base_ptrEv(ptr noundef nonnull align 8 dereferenceable(32) %5) #10
  store ptr %call5, ptr %__base_z, align 8
  %6 = load i8, ptr %__insert_left, align 1
  %loadedv = trunc i8 %6 to i1
  %7 = load ptr, ptr %__base_z, align 8
  %8 = load ptr, ptr %__p.addr, align 8
  %_M_impl = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr = getelementptr inbounds i8, ptr %_M_impl, i64 8
  %_M_header = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr, i32 0, i32 0
  call void @_ZNSt9__rb_tree12_Node_traitsISt4pairIKPK1VjEPS6_E23_S_insert_and_rebalanceEbPSt18_Rb_tree_node_baseSA_RS9_(i1 noundef zeroext %loadedv, ptr noundef %7, ptr noundef %8, ptr noundef nonnull align 8 dereferenceable(32) %_M_header) #10
  %_M_impl6 = getelementptr inbounds nuw %"class.std::_Rb_tree", ptr %this1, i32 0, i32 0
  %add.ptr7 = getelementptr inbounds i8, ptr %_M_impl6, i64 8
  %_M_node_count = getelementptr inbounds nuw %"struct.std::_Rb_tree_header", ptr %add.ptr7, i32 0, i32 1
  %9 = load i64, ptr %_M_node_count, align 8
  %inc = add i64 %9, 1
  store i64 %inc, ptr %_M_node_count, align 8
  %10 = load ptr, ptr %__base_z, align 8
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %retval, ptr noundef %10) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %11 = load ptr, ptr %coerce.dive, align 8
  ret ptr %11
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt9__rb_tree12_Node_traitsISt4pairIKPK1VjEPS6_E23_S_insert_and_rebalanceEbPSt18_Rb_tree_node_baseSA_RS9_(i1 noundef zeroext %__insert_left, ptr noundef nonnull %__x, ptr noundef nonnull %__p, ptr noundef nonnull align 8 dereferenceable(32) %__header) #1 comdat align 2 {
entry:
  %__insert_left.addr = alloca i8, align 1
  %__x.addr = alloca ptr, align 8
  %__p.addr = alloca ptr, align 8
  %__header.addr = alloca ptr, align 8
  %storedv = zext i1 %__insert_left to i8
  store i8 %storedv, ptr %__insert_left.addr, align 1
  store ptr %__x, ptr %__x.addr, align 8
  store ptr %__p, ptr %__p.addr, align 8
  store ptr %__header, ptr %__header.addr, align 8
  %0 = load i8, ptr %__insert_left.addr, align 1
  %loadedv = trunc i8 %0 to i1
  %1 = load ptr, ptr %__x.addr, align 8
  %2 = load ptr, ptr %__p.addr, align 8
  %3 = load ptr, ptr %__header.addr, align 8, !nonnull !10, !align !11
  call void @_ZSt29_Rb_tree_insert_and_rebalancebPSt18_Rb_tree_node_baseS0_RS_(i1 noundef zeroext %loadedv, ptr noundef %1, ptr noundef %2, ptr noundef nonnull align 8 dereferenceable(32) %3) #10
  ret void
}

; Function Attrs: nounwind
declare void @_ZSt29_Rb_tree_insert_and_rebalancebPSt18_Rb_tree_node_baseS0_RS_(i1 noundef zeroext, ptr noundef, ptr noundef, ptr noundef nonnull align 8 dereferenceable(32)) #8

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt11_Tuple_implILm0EJRKPK1VEEC2ES4_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(8) %__head) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__head.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__head, ptr %__head.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %0 = load ptr, ptr %__head.addr, align 8, !nonnull !10, !align !11
  call void @_ZNSt10_Head_baseILm0ERKPK1VLb0EEC2ES4_(ptr noundef nonnull align 8 dereferenceable(8) %this1, ptr noundef nonnull align 8 dereferenceable(8) %0)
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local void @_ZNSt10_Head_baseILm0ERKPK1VLb0EEC2ES4_(ptr noundef nonnull align 8 dereferenceable(8) %this, ptr noundef nonnull align 8 dereferenceable(8) %__h) unnamed_addr #1 comdat align 2 {
entry:
  %this.addr = alloca ptr, align 8
  %__h.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  store ptr %__h, ptr %__h.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %_M_head_impl = getelementptr inbounds nuw %"struct.std::_Head_base", ptr %this1, i32 0, i32 0
  %0 = load ptr, ptr %__h.addr, align 8, !nonnull !10, !align !11
  store ptr %0, ptr %_M_head_impl, align 8
  ret void
}

; Function Attrs: mustprogress noinline nounwind optnone uwtable
define linkonce_odr dso_local ptr @_ZNSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE3endEv(ptr noundef nonnull align 8 dereferenceable(48) %this) #1 comdat align 2 {
entry:
  %retval = alloca %"struct.std::_Rb_tree_iterator", align 8
  %this.addr = alloca ptr, align 8
  store ptr %this, ptr %this.addr, align 8
  %this1 = load ptr, ptr %this.addr, align 8
  %call = call noundef ptr @_ZNKSt8_Rb_treeIPK1VSt4pairIKS2_jESt10_Select1stIS5_ESt4lessIS2_ESaIS5_EE6_M_endEv(ptr noundef nonnull align 8 dereferenceable(48) %this1) #10
  call void @_ZNSt17_Rb_tree_iteratorISt4pairIKPK1VjEEC2EPSt18_Rb_tree_node_base(ptr noundef nonnull align 8 dereferenceable(8) %retval, ptr noundef %call) #10
  %coerce.dive = getelementptr inbounds nuw %"struct.std::_Rb_tree_iterator", ptr %retval, i32 0, i32 0
  %0 = load ptr, ptr %coerce.dive, align 8
  ret ptr %0
}

; Function Attrs: nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: readwrite)
declare void @llvm.experimental.noalias.scope.decl(metadata) #9

attributes #0 = { mustprogress noinline norecurse nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #1 = { mustprogress noinline nounwind optnone uwtable "frame-pointer"="all" "min-legal-vector-width"="0" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #2 = { "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #3 = { nocallback nofree nounwind willreturn memory(argmem: readwrite) }
attributes #4 = { nounwind willreturn memory(read) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #5 = { noreturn "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #6 = { nobuiltin allocsize(0) "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #7 = { nobuiltin nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #8 = { nounwind "frame-pointer"="all" "no-trapping-math"="true" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+cmov,+cx8,+fxsr,+mmx,+sse,+sse2,+x87" "tune-cpu"="generic" }
attributes #9 = { nocallback nofree nosync nounwind willreturn memory(inaccessiblemem: readwrite) }
attributes #10 = { nounwind }
attributes #11 = { nounwind willreturn memory(read) }
attributes #12 = { noreturn }
attributes #13 = { builtin allocsize(0) }
attributes #14 = { builtin nounwind }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{i32 8, !"PIC Level", i32 2}
!2 = !{i32 7, !"PIE Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 2}
!4 = !{i32 7, !"frame-pointer", i32 2}
!5 = !{!"clang version 22.1.8 (https://github.com/conda-forge/clangdev-feedstock 015bdba1263c0b3ebb3c518ff5947fbd99692bd0)"}
!6 = distinct !{!6, !7}
!7 = !{!"llvm.loop.mustprogress"}
!8 = distinct !{!8, !7}
!9 = distinct !{!9, !7}
!10 = !{}
!11 = !{i64 8}
!12 = distinct !{!12, !7}
!13 = !{!14}
!14 = distinct !{!14, !15, !"_ZN9__gnu_cxx14__alloc_traitsISaISt13_Rb_tree_nodeISt4pairIKPK1VjEEES8_E17_S_select_on_copyERKS9_: %agg.result"}
!15 = distinct !{!15, !"_ZN9__gnu_cxx14__alloc_traitsISaISt13_Rb_tree_nodeISt4pairIKPK1VjEEES8_E17_S_select_on_copyERKS9_"}
!16 = !{!17}
!17 = distinct !{!17, !18, !"_ZNSt16allocator_traitsISaISt13_Rb_tree_nodeISt4pairIKPK1VjEEEE37select_on_container_copy_constructionERKS8_: %agg.result"}
!18 = distinct !{!18, !"_ZNSt16allocator_traitsISaISt13_Rb_tree_nodeISt4pairIKPK1VjEEEE37select_on_container_copy_constructionERKS8_"}
!19 = distinct !{!19, !7}
!20 = distinct !{!20, !7}
!21 = distinct !{!21, !7}
!22 = distinct !{!22, !7}
!23 = distinct !{!23, !7}
!24 = distinct !{!24, !7}
