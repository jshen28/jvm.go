package lang

import (
	. "jvmgo/any"
	cp "jvmgo/classpath"
	"jvmgo/jvm/rtda"
	rtc "jvmgo/jvm/rtda/class"
	"strings"
)

func init() {
	_class(getClassLoader0, "getClassLoader0", "()Ljava/lang/ClassLoader;")
	_class(getComponentType, "getComponentType", "()Ljava/lang/Class;")
	_class(getConstantPool, "getConstantPool", "()Lsun/reflect/ConstantPool;")
	_class(getDeclaringClass0, "getDeclaringClass0", "()Ljava/lang/Class;")
	_class(getEnclosingMethod0, "getEnclosingMethod0", "()[Ljava/lang/Object;")
	_class(getInterfaces0, "getInterfaces0", "()[Ljava/lang/Class;")
	_class(getModifiers, "getModifiers", "()I")
	_class(getName0, "getName0", "()Ljava/lang/String;")
	_class(getSuperclass, "getSuperclass", "()Ljava/lang/Class;")
	_class(isArray, "isArray", "()Z")
	_class(isAssignableFrom, "isAssignableFrom", "(Ljava/lang/Class;)Z")
	_class(isInstance, "isInstance", "(Ljava/lang/Object;)Z")
	_class(isInterface, "isInterface", "()Z")
	_class(isPrimitive, "isPrimitive", "()Z")
}

func _class(method Any, name, desc string) {
	rtc.RegisterNativeMethod("java/lang/Class", name, desc, method)
}

// native ClassLoader getClassLoader0();
// ()Ljava/lang/ClassLoader;
func getClassLoader0(frame *rtda.Frame) {
	class := _popClass(frame)
	from := class.LoadedFrom()

	stack := frame.OperandStack()
	if cp.IsBootClassPath(from) {
		stack.PushRef(nil)
		return
	}

	clClass := class.ClassLoader().LoadClass("java/lang/ClassLoader")
	getSysCl := clClass.GetStaticMethod("getSystemClassLoader", "()Ljava/lang/ClassLoader;")
	frame.Thread().InvokeMethod(getSysCl)
}

// public native Class<?> getComponentType();
// ()Ljava/lang/Class;
func getComponentType(frame *rtda.Frame) {
	class := _popClass(frame)
	componentClass := class.ComponentClass()
	componentClassObj := componentClass.JClass()

	stack := frame.OperandStack()
	stack.PushRef(componentClassObj)
}

// native ConstantPool getConstantPool();
// ()Lsun/reflect/ConstantPool;
func getConstantPool(frame *rtda.Frame) {
	class := _popClass(frame)
	cpClass := class.ClassLoader().LoadClass("sun/reflect/ConstantPool")
	if cpClass.InitializationNotStarted() {
		frame.RevertNextPC()
		frame.Thread().InitClass(cpClass)
		return
	}

	cp := class.ConstantPool()
	cpObj := cpClass.NewObjWithExtra(cp) // todo init cpObj
	frame.OperandStack().PushRef(cpObj)
}

// private native Class<?> getDeclaringClass0();
// ()Ljava/lang/Class;
func getDeclaringClass0(frame *rtda.Frame) {
	class := _popClass(frame)
	if class.IsArray() || class.IsPrimitive() {
		frame.OperandStack().PushRef(nil)
		return
	}

	lastDollerIndex := strings.LastIndex(class.Name(), "$")
	if lastDollerIndex < 0 {
		frame.OperandStack().PushRef(nil)
		return
	}

	// todo
	declaringClassName := class.Name()[:lastDollerIndex]
	declaringClass := frame.ClassLoader().LoadClass(declaringClassName)
	frame.OperandStack().PushRef(declaringClass.JClass())
}

// private native Object[] getEnclosingMethod0();
// ()[Ljava/lang/Object;
func getEnclosingMethod0(frame *rtda.Frame) {
	class := _popClass(frame)
	emInfo := class.Attributes().EnclosingMethod()
	emInfoObj := _createEnclosintMethodInfo(frame.ClassLoader(), emInfo)
	frame.OperandStack().PushRef(emInfoObj)
}

func _createEnclosintMethodInfo(classLoader *rtc.ClassLoader, emInfo *rtc.EnclosingMethod) *rtc.Obj {
	if emInfo == nil {
		return nil
	}

	enclosingClass := classLoader.LoadClass(emInfo.ClassName())
	enclosingClassObj := enclosingClass.JClass()
	var methodNameObj, methodDescriptorObj *rtc.Obj
	if emInfo.MethodName() != "" {
		methodNameObj = rtda.NewJString(emInfo.MethodName(), classLoader)
		methodDescriptorObj = rtda.NewJString(emInfo.MethodDescriptor(), classLoader)
	} else {
		methodNameObj, methodDescriptorObj = nil, nil
	}

	objs := []*rtc.Obj{enclosingClassObj, methodNameObj, methodDescriptorObj}
	return rtc.NewRefArray2(classLoader.JLObjectClass(), objs) // Object[]
}

// private native Class<?>[] getInterfaces0();
// ()[Ljava/lang/Class;
func getInterfaces0(frame *rtda.Frame) {
	class := _popClass(frame)
	interfaces := class.Interfaces()
	interfaceObjs := make([]*rtc.Obj, len(interfaces))
	for i, iface := range interfaces {
		interfaceObjs[i] = iface.JClass()
	}

	jlClassClass := class.ClassLoader().JLClassClass()
	interfaceArr := rtc.NewRefArray2(jlClassClass, interfaceObjs)

	stack := frame.OperandStack()
	stack.PushRef(interfaceArr)
}

// private native String getName0();
// ()Ljava/lang/String;
func getName0(frame *rtda.Frame) {
	class := _popClass(frame)
	name := class.NameJlsFormat()
	nameObj := rtda.NewJString(name, frame)

	stack := frame.OperandStack()
	stack.PushRef(nameObj)
}

// public native int getModifiers();
// ()I
func getModifiers(frame *rtda.Frame) {
	class := _popClass(frame)
	modifiers := class.GetAccessFlags()

	stack := frame.OperandStack()
	stack.PushInt(int32(modifiers))
}

// public native Class<? super T> getSuperclass();
// ()Ljava/lang/Class;
func getSuperclass(frame *rtda.Frame) {
	class := _popClass(frame)
	superClass := class.SuperClass()

	stack := frame.OperandStack()
	if superClass != nil {
		stack.PushRef(superClass.JClass())
	} else {
		stack.PushNull()
	}
}

// public native boolean isAssignableFrom(Class<?> cls);
// (Ljava/lang/Class;)Z
func isAssignableFrom(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()
	cls := vars.GetRef(1)

	thisClass := this.Extra().(*rtc.Class)
	clsClass := cls.Extra().(*rtc.Class)
	ok := thisClass.IsAssignableFrom(clsClass)

	stack := frame.OperandStack()
	stack.PushBoolean(ok)
}

// public native boolean isInstance(Object obj);
// (Ljava/lang/Object;)Z
func isInstance(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()
	obj := vars.GetRef(1)

	class := this.Extra().(*rtc.Class)
	ok := obj.IsInstanceOf(class)

	stack := frame.OperandStack()
	stack.PushBoolean(ok)
}

// public native boolean isArray();
// ()Z
func isArray(frame *rtda.Frame) {
	class := _popClass(frame)
	stack := frame.OperandStack()
	stack.PushBoolean(class.IsArray())
}

// public native boolean isInterface();
// ()Z
func isInterface(frame *rtda.Frame) {
	class := _popClass(frame)
	stack := frame.OperandStack()
	stack.PushBoolean(class.IsInterface())
}

// public native boolean isPrimitive();
// ()Z
func isPrimitive(frame *rtda.Frame) {
	class := _popClass(frame)
	stack := frame.OperandStack()
	stack.PushBoolean(class.IsPrimitive())
}

func _popClass(frame *rtda.Frame) *rtc.Class {
	vars := frame.LocalVars()
	this := vars.GetThis()
	return this.Extra().(*rtc.Class)
}
