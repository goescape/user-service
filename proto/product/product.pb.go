// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: proto/product/product.proto

package product

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProductInsertRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Price         float32                `protobuf:"fixed32,3,opt,name=price,proto3" json:"price,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Qty           uint32                 `protobuf:"varint,5,opt,name=qty,proto3" json:"qty,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProductInsertRequest) Reset() {
	*x = ProductInsertRequest{}
	mi := &file_proto_product_product_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductInsertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductInsertRequest) ProtoMessage() {}

func (x *ProductInsertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductInsertRequest.ProtoReflect.Descriptor instead.
func (*ProductInsertRequest) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{0}
}

func (x *ProductInsertRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ProductInsertRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProductInsertRequest) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ProductInsertRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ProductInsertRequest) GetQty() uint32 {
	if x != nil {
		return x.Qty
	}
	return 0
}

type ProductInsertResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           string                 `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProductInsertResponse) Reset() {
	*x = ProductInsertResponse{}
	mi := &file_proto_product_product_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductInsertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductInsertResponse) ProtoMessage() {}

func (x *ProductInsertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductInsertResponse.ProtoReflect.Descriptor instead.
func (*ProductInsertResponse) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{1}
}

func (x *ProductInsertResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type Product struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Price         float32                `protobuf:"fixed32,4,opt,name=price,proto3" json:"price,omitempty"`
	Description   string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Qty           uint32                 `protobuf:"varint,6,opt,name=qty,proto3" json:"qty,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Product) Reset() {
	*x = Product{}
	mi := &file_proto_product_product_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{2}
}

func (x *Product) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Product) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Product) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Product) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Product) GetQty() uint32 {
	if x != nil {
		return x.Qty
	}
	return 0
}

type ListProductRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          uint32                 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit         uint32                 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	ProductIds    string                 `protobuf:"bytes,3,opt,name=product_ids,json=productIds,proto3" json:"product_ids,omitempty"` // comma separated
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProductRequest) Reset() {
	*x = ListProductRequest{}
	mi := &file_proto_product_product_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProductRequest) ProtoMessage() {}

func (x *ListProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProductRequest.ProtoReflect.Descriptor instead.
func (*ListProductRequest) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{3}
}

func (x *ListProductRequest) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListProductRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListProductRequest) GetProductIds() string {
	if x != nil {
		return x.ProductIds
	}
	return ""
}

type Meta struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TotalData     uint32                 `protobuf:"varint,1,opt,name=total_data,json=totalData,proto3" json:"total_data,omitempty"`
	TotalPage     uint32                 `protobuf:"varint,2,opt,name=total_page,json=totalPage,proto3" json:"total_page,omitempty"`
	CurrentPage   uint32                 `protobuf:"varint,3,opt,name=current_page,json=currentPage,proto3" json:"current_page,omitempty"`
	Limit         uint32                 `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Meta) Reset() {
	*x = Meta{}
	mi := &file_proto_product_product_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Meta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meta) ProtoMessage() {}

func (x *Meta) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Meta.ProtoReflect.Descriptor instead.
func (*Meta) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{4}
}

func (x *Meta) GetTotalData() uint32 {
	if x != nil {
		return x.TotalData
	}
	return 0
}

func (x *Meta) GetTotalPage() uint32 {
	if x != nil {
		return x.TotalPage
	}
	return 0
}

func (x *Meta) GetCurrentPage() uint32 {
	if x != nil {
		return x.CurrentPage
	}
	return 0
}

func (x *Meta) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type ListProductResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*Product             `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Meta          *Meta                  `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProductResponse) Reset() {
	*x = ListProductResponse{}
	mi := &file_proto_product_product_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProductResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProductResponse) ProtoMessage() {}

func (x *ListProductResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProductResponse.ProtoReflect.Descriptor instead.
func (*ListProductResponse) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{5}
}

func (x *ListProductResponse) GetItems() []*Product {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ListProductResponse) GetMeta() *Meta {
	if x != nil {
		return x.Meta
	}
	return nil
}

type ReduceProductsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*ProductItem         `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReduceProductsRequest) Reset() {
	*x = ReduceProductsRequest{}
	mi := &file_proto_product_product_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReduceProductsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReduceProductsRequest) ProtoMessage() {}

func (x *ReduceProductsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReduceProductsRequest.ProtoReflect.Descriptor instead.
func (*ReduceProductsRequest) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{6}
}

func (x *ReduceProductsRequest) GetItems() []*ProductItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type ProductItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     string                 `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Qty           uint32                 `protobuf:"varint,2,opt,name=qty,proto3" json:"qty,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProductItem) Reset() {
	*x = ProductItem{}
	mi := &file_proto_product_product_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductItem) ProtoMessage() {}

func (x *ProductItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductItem.ProtoReflect.Descriptor instead.
func (*ProductItem) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{7}
}

func (x *ProductItem) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *ProductItem) GetQty() uint32 {
	if x != nil {
		return x.Qty
	}
	return 0
}

type ReduceProductsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           string                 `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReduceProductsResponse) Reset() {
	*x = ReduceProductsResponse{}
	mi := &file_proto_product_product_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReduceProductsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReduceProductsResponse) ProtoMessage() {}

func (x *ReduceProductsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_product_product_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReduceProductsResponse.ProtoReflect.Descriptor instead.
func (*ReduceProductsResponse) Descriptor() ([]byte, []int) {
	return file_proto_product_product_proto_rawDescGZIP(), []int{8}
}

func (x *ReduceProductsResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_proto_product_product_proto protoreflect.FileDescriptor

const file_proto_product_product_proto_rawDesc = "" +
	"\n" +
	"\x1bproto/product/product.proto\x12\x05proto\"\x8d\x01\n" +
	"\x14ProductInsertRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x14\n" +
	"\x05price\x18\x03 \x01(\x02R\x05price\x12 \n" +
	"\vdescription\x18\x04 \x01(\tR\vdescription\x12\x10\n" +
	"\x03qty\x18\x05 \x01(\rR\x03qty\")\n" +
	"\x15ProductInsertResponse\x12\x10\n" +
	"\x03msg\x18\x01 \x01(\tR\x03msg\"\x90\x01\n" +
	"\aProduct\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12\x14\n" +
	"\x05price\x18\x04 \x01(\x02R\x05price\x12 \n" +
	"\vdescription\x18\x05 \x01(\tR\vdescription\x12\x10\n" +
	"\x03qty\x18\x06 \x01(\rR\x03qty\"_\n" +
	"\x12ListProductRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\rR\x04page\x12\x14\n" +
	"\x05limit\x18\x02 \x01(\rR\x05limit\x12\x1f\n" +
	"\vproduct_ids\x18\x03 \x01(\tR\n" +
	"productIds\"}\n" +
	"\x04Meta\x12\x1d\n" +
	"\n" +
	"total_data\x18\x01 \x01(\rR\ttotalData\x12\x1d\n" +
	"\n" +
	"total_page\x18\x02 \x01(\rR\ttotalPage\x12!\n" +
	"\fcurrent_page\x18\x03 \x01(\rR\vcurrentPage\x12\x14\n" +
	"\x05limit\x18\x04 \x01(\rR\x05limit\"\\\n" +
	"\x13ListProductResponse\x12$\n" +
	"\x05items\x18\x01 \x03(\v2\x0e.proto.ProductR\x05items\x12\x1f\n" +
	"\x04meta\x18\x02 \x01(\v2\v.proto.MetaR\x04meta\"A\n" +
	"\x15ReduceProductsRequest\x12(\n" +
	"\x05items\x18\x01 \x03(\v2\x12.proto.ProductItemR\x05items\">\n" +
	"\vProductItem\x12\x1d\n" +
	"\n" +
	"product_id\x18\x01 \x01(\tR\tproductId\x12\x10\n" +
	"\x03qty\x18\x02 \x01(\rR\x03qty\"*\n" +
	"\x16ReduceProductsResponse\x12\x10\n" +
	"\x03msg\x18\x01 \x01(\tR\x03msg2\xf1\x01\n" +
	"\x0eProductService\x12J\n" +
	"\rInsertProduct\x12\x1b.proto.ProductInsertRequest\x1a\x1c.proto.ProductInsertResponse\x12D\n" +
	"\vListProduct\x12\x19.proto.ListProductRequest\x1a\x1a.proto.ListProductResponse\x12M\n" +
	"\x0eReduceProducts\x12\x1c.proto.ReduceProductsRequest\x1a\x1d.proto.ReduceProductsResponseB\x0fZ\rproto/productb\x06proto3"

var (
	file_proto_product_product_proto_rawDescOnce sync.Once
	file_proto_product_product_proto_rawDescData []byte
)

func file_proto_product_product_proto_rawDescGZIP() []byte {
	file_proto_product_product_proto_rawDescOnce.Do(func() {
		file_proto_product_product_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_product_product_proto_rawDesc), len(file_proto_product_product_proto_rawDesc)))
	})
	return file_proto_product_product_proto_rawDescData
}

var file_proto_product_product_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_product_product_proto_goTypes = []any{
	(*ProductInsertRequest)(nil),   // 0: proto.ProductInsertRequest
	(*ProductInsertResponse)(nil),  // 1: proto.ProductInsertResponse
	(*Product)(nil),                // 2: proto.Product
	(*ListProductRequest)(nil),     // 3: proto.ListProductRequest
	(*Meta)(nil),                   // 4: proto.Meta
	(*ListProductResponse)(nil),    // 5: proto.ListProductResponse
	(*ReduceProductsRequest)(nil),  // 6: proto.ReduceProductsRequest
	(*ProductItem)(nil),            // 7: proto.ProductItem
	(*ReduceProductsResponse)(nil), // 8: proto.ReduceProductsResponse
}
var file_proto_product_product_proto_depIdxs = []int32{
	2, // 0: proto.ListProductResponse.items:type_name -> proto.Product
	4, // 1: proto.ListProductResponse.meta:type_name -> proto.Meta
	7, // 2: proto.ReduceProductsRequest.items:type_name -> proto.ProductItem
	0, // 3: proto.ProductService.InsertProduct:input_type -> proto.ProductInsertRequest
	3, // 4: proto.ProductService.ListProduct:input_type -> proto.ListProductRequest
	6, // 5: proto.ProductService.ReduceProducts:input_type -> proto.ReduceProductsRequest
	1, // 6: proto.ProductService.InsertProduct:output_type -> proto.ProductInsertResponse
	5, // 7: proto.ProductService.ListProduct:output_type -> proto.ListProductResponse
	8, // 8: proto.ProductService.ReduceProducts:output_type -> proto.ReduceProductsResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_product_product_proto_init() }
func file_proto_product_product_proto_init() {
	if File_proto_product_product_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_product_product_proto_rawDesc), len(file_proto_product_product_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_product_product_proto_goTypes,
		DependencyIndexes: file_proto_product_product_proto_depIdxs,
		MessageInfos:      file_proto_product_product_proto_msgTypes,
	}.Build()
	File_proto_product_product_proto = out.File
	file_proto_product_product_proto_goTypes = nil
	file_proto_product_product_proto_depIdxs = nil
}
