package cobbles

import (
    "reflect"
)

// Takes the data from the given "in" interface and places it the given "out" interface.
func Cast(in interface{}, out interface{}) {
    cast := &caster{}
    cast.marshal(reflect.ValueOf(in), reflect.ValueOf(out))
}

type caster struct {

}

// For each key/type in "to", attempt to find a source in "from".
func (this *caster) marshal(in reflect.Value, out reflect.Value) {
    switch out.Kind() {
    case reflect.Interface:
        if out.IsNil() {
            this.nilv(in, out)
        } else {
            this.marshal(in, out.Elem())
        }
    case reflect.Map:
        this.mapv(out, in)
    case reflect.Ptr:
        if out.IsNil() {
            this.nilv(in, out)
        } else {
            this.marshal(in, out.Elem())
        }
    case reflect.Struct:
        this.structv(in, out)
    case reflect.Slice:
        this.slicev(in, out)
    case reflect.String:
        this.stringv(in, out)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        this.intv(in, out)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        this.uintv(in, out)
    case reflect.Float32, reflect.Float64:
        this.floatv(in, out)
    case reflect.Bool:
        this.boolv(in, out)
    default:
        panic("Can't marshal type yet: " + out.Type().String())
    }
}

func (this *caster) nilv(in reflect.Value, out reflect.Value) {

}

func (this *caster) mapv(in reflect.Value, out reflect.Value) {
    
}

func (this *caster) structv(in reflect.Value, out reflect.Value) {
    
}

func (this *caster) slicev(in reflect.Value, out reflect.Value) {
    
}

func (this *caster) stringv(in reflect.Value, out reflect.Value) {
    
}

func (this *caster) intv(in reflect.Value, out reflect.Value) {
    
}

func (this *caster) uintv(in reflect.Value, out reflect.Value) {
    
}

func (this *caster) floatv(in reflect.Value, out reflect.Value) {
    
}

func (this *caster) boolv(in reflect.Value, out reflect.Value) {
    
}