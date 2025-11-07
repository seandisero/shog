package shog

import (
	"bytes"
	"strconv"
	"strings"
)

type Image struct {
	size   UV
	origin UV
	Data   []rune
}

var TEST_IMAGE = Image{
	size:   NewUV(2, 2),
	origin: NewUV(32, 16),
	Data:   []rune("/\\\\/"),
}

var TEST_IMAGE2 = Image{
	size:   NewUV(16, 16),
	origin: NewUV(38, 16),
	Data: []rune(strings.ReplaceAll(`
   ##########   
  ##        ##  
 ##          ## 
 ##  ##  ##  ## 
 ##  ##  ##  ## 
 ##          ## 
 ##          ## 
 ##  ##  ##  ## 
 ## ##    ## ## 
  ##  ####  ##  
   ##########   `, "\n", "")),
}

var image3 string = `24:24
*       .        *      
  .    *      .      *  
       .   *      .     
 *   .:::::::::..    .  
    .::::::::::::::. *  
  * :::::::::::::::::.  
 .  ::::::'~~~'::::::: *
   *:::::'~~~~~'::::::. 
 .  ::::~~~~~~~~~::::  .
*   :::~~~~~~~~~~~~::: *
   .:::~~~~~~~~~~~~:::. 
 *  :::~~~~~~~~~~~~::: .
    .:::~~~~~~~~~~:::*  
  .  ::::'~~~~~'::::  . 
 *   '::::'~~~'::::'  * 
   .  '::::::::::'   *  
  *   . ':::::::'  .    
 .  *     ':::'      *  
    .  *    '    .    . 
  *      .    *      *  
 .   *      .     *   . 
   .    *      .    *   
 *   .      *     .   * 
  .    *       .      . `

var TEST_IMAGE3 *Image = ImageFromBytes([]byte(image3))

func (i *Image) SetOrigin(uv UV) {
	i.origin = uv
}

func ImageFromBytes(b []byte) *Image {
	split := bytes.Split(b, []byte("\n"))
	if len(split) < 2 {
		return nil
	}
	sizeBytes := bytes.Split(split[0], []byte(":"))
	w, err := strconv.Atoi(string(sizeBytes[0]))
	h, err := strconv.Atoi(string(sizeBytes[0]))
	if err != nil {
		return nil
	}
	image := Image{
		size: NewUV(w, h),
		Data: []rune(string(bytes.ReplaceAll(bytes.Join(split[1:], []byte("")), []byte("\n"), []byte("")))),
	}

	return &image
}
