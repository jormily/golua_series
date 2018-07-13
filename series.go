package series

import(
	"github.com/yuin/gopher-lua"
	"bytes"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"encode":        encode,
		"decode":        decode,
	})
	L.Push(mod)
	return 1
}


/*
	    table_key      table_value
	[int8 int32 数据][int8 int32 数据]
*/
func encodeTable(L *lua.LState,tb *lua.LTable) bytes.Buffer  {
	var body bytes.Buffer 
	tb.ForEach(func (key lua.LValue,val lua.LValue)  {
		switch key.Type() {
			case lua.LTNumber:
				dtype := byte(lua.LTNumber)
				data := lua.LVAsNumber(key)

				body.WriteByte(byte(dtype))
				body.Write(Float64ToByte(float64(data)))
			case lua.LTString:
				dtype := byte(lua.LTString)
				data := lua.LVAsString(key)
				length := len(data)

				body.WriteByte(byte(dtype))
				body.Write(Float32ToByte(float32(length)))
				body.WriteString(string(data))
			default:
				goto Cend
		}

		switch val.Type() {
			case lua.LTNumber:
				dtype := byte(lua.LTNumber)
				data := lua.LVAsNumber(val)

				body.WriteByte(byte(dtype))
				body.Write(Float64ToByte(float64(data)))
			case lua.LTString:
				dtype := byte(lua.LTString)
				data := lua.LVAsString(val)
				length := len(data)

				body.WriteByte(byte(dtype))
				body.Write(Float32ToByte(float32(length)))
				body.WriteString(string(data))
			case lua.LTBool:
				dtype := byte(lua.LTBool)
				body.WriteByte(byte(dtype))
				if lua.LVAsBool(val) {
					body.WriteByte(byte(1))
				}else{
					body.WriteByte(byte(0))
				}
			case lua.LTUserData:
				dtype := byte(lua.LTUserData)
				body.WriteByte(byte(dtype))
				userData := val.(*lua.LUserData)
				buff := userData.Value.(*bytes.Buffer)
				body.Write(Float32ToByte(float32(buff.Len())))
				body.Write(buff.Bytes())
			case lua.LTTable:
				dtype := byte(lua.LTTable)
				body.WriteByte(byte(dtype))
				intb := val.(*lua.LTable)
				buff := encodeTable(L,intb)
				body.Write(Float32ToByte(float32(buff.Len())))
				body.Write(buff.Bytes())
			default:			

		}		
		Cend:
	})

	return body
}

func encode(L *lua.LState) int  {
	tb := L.CheckTable(1)
	if tb == nil { return 0 }
	buff := encodeTable(L,tb)
	if buff.Len() > 0 {
		userData := L.NewUserData()
		userData.Value = buff
		L.Push(userData)
		return 1
	}

	return 0
}

func decdeTable(L *lua.LState,buff *bytes.Buffer) *lua.LTable  {
	var tb *lua.LTable = L.NewTable()
	var lkey lua.LValue
	var lval lua.LValue 
	var flag bool

	for ;buff.Len() > 0; {
		kt,ok := buff.ReadByte()
		if ok != nil { return nil }
		flag = true

		switch lua.LValueType(kt) {
			case lua.LTNumber:
				bts := buff.Next(8)
				if len(bts) < 8 { goto wrong }
				lkey = lua.LNumber(ByteToFloat64(bts))
			case lua.LTString:
				bts := buff.Next(4)
				if len(bts) < 4 { goto wrong }
				length := int(ByteToFloat32(bts))
				bts = buff.Next(length)
				if len(bts) < length { goto wrong }
				lkey = lua.LValue(lua.LString(string(bts)))
			default:
				flag = false
		} 

		kt,ok = buff.ReadByte()
		if ok != nil { return nil }

		switch lua.LValueType(kt) {
			case lua.LTNumber:
				bts := buff.Next(8)
				if len(bts) < 8 { goto wrong }
				lval = lua.LNumber(ByteToFloat64(bts))
			case lua.LTString:
				bts := buff.Next(4)
				if len(bts) < 4 { goto wrong }
				length := int(ByteToFloat32(bts))
				bts = buff.Next(length)
				if len(bts) < length { goto wrong }
				lval = lua.LValue(lua.LString(string(bts)))
			case lua.LTBool:
				bts := buff.Next(1)
				if len(bts) < 1 { goto wrong}
				lval = lua.LBool(bts[0] == 1)
			case lua.LTUserData:
				bts := buff.Next(4)
				if len(bts) < 4 { goto wrong }
				length := int(ByteToFloat32(bts))
				bts = buff.Next(length)
				if len(bts) < length { goto wrong }
				
				userData := L.NewUserData()
				userData.Value = bytes.NewBuffer(bts)
				lval = userData
			case lua.LTTable:
				bts := buff.Next(4)
				if len(bts) < 4 { goto wrong }
				length := int(ByteToFloat32(bts))
				bts = buff.Next(length)
				if len(bts) < length { goto wrong }

				intb := decdeTable(L,bytes.NewBuffer(bts))
				if intb == nil { goto wrong }
				lval = intb
			default:
				flag = false
		} 

		if flag {
			tb.RawSet(lkey,lval)
		}
	}
		
	return tb
	
	wrong:
	return nil
}

func decode(L *lua.LState) int  {
	userData := L.CheckUserData(1)
	if userData == nil { return 0 }
	buff := userData.Value.(bytes.Buffer)
	tb := decdeTable(L,&buff)
	if tb == nil { return 0 }

	L.Push(tb)
	return 1
}