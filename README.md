# symwalk

```
import "github.com/edwardrf/symwalk"
```
symwalk provide a walk function similar to filepath.Walk but follow symbolic links
but avoid visiting directories more than once even if there is a symbolic link loop.

### Whats the difference from [facebookarchive/symwalk](https://github.com/facebookarchive/symwalk)
The main difference is this Walk implementation will not visit the same directories more
than once even if there is a symbolic link loop.
This is achieved by 
 - always evaluate the symbolic links to its real path by filepath.EvalSymlinks, this means your walk function will received the real path instead of the relative path from the symbolic links
 - keep record of all visited directories using a map

Since the underlying implementation make use of filepath.Walk, the rest of the behavior is the same as filepath.Walk.

### Example
```
symwalk.Walk("/tmp/", func(path string, info os.FileInfo, err error) error {
  fmt.Println(path)
  return nil
}) 
```
