//go:build ignore

package main

var ScreenSize vec2
var TextureColors [1]vec3 // TODO does this name actually make sense?
var SourcePalette [1]vec3
var Palette0 [1]vec3

// TODO okay so this is a bit awkward because I have to hard code the array
// size for the above arrays. The problem with that is I want this to be
// generic, but I don't know how big my palettes will be in advance. One thing
// I could try is to just include a bunch of buffer and have some uniform var
// for the effective array size. I think I'm better off sticking with Palette0,
// Palette1, etc being separate variables rather than turning it into a 2d
// array because it would be really annoying to fill those in with them being
// flattened. Maybe I'm wrong about that IDK. It *is* pretty ugly I will admit.

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	var out vec3

    out.rgb = imageSrc0UnsafeAt(texCoord).rgb
    
    for i := 0; i < len(TextureColors); i++ {
        if (imageSrc1UnsafeAt(texCoord).rgb * 255.0) == TextureColors[i] {
            for j := 0; j < len(SourcePalette); j++ {
                if (imageSrc0UnsafeAt(texCoord).rgb * 255.0) == SourcePalette[j] {
                    if i == 0 { // TODO ugly but necessary?
                        out.rgb = Palette0[j].rgb
                    }
                }
            }
        }
    }

	return vec4(out, 1.0)
}
