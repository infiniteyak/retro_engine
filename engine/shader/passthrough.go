//go:build ignore

package main

var ScreenSize vec2

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
    /*
	center := ScreenSize / 2
	amount := (center - Cursor) / 10 / imageSrcTextureSize()
	var clr vec3
	clr.r = imageSrc2At(texCoord + amount).r
	clr.g = imageSrc2UnsafeAt(texCoord).g
	clr.b = imageSrc2At(texCoord - amount).b
    */
	var clr vec3
	clr.r = imageSrc0UnsafeAt(texCoord).r
	clr.g = imageSrc0UnsafeAt(texCoord).g
	clr.b = imageSrc0UnsafeAt(texCoord).b

	return vec4(clr, 1.0)
}
