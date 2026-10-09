package launch

import (
 "errors"
 "flag"
 "os"
 "path/filepath"
 "testing"
)
func TestProjectPath(t *testing.T) {
 root:=t.TempDir()
 noProject,err:=ProjectPath(nil)
 if err!=nil||noProject!=""{t.Fatalf("empty args: %q %v",noProject,err)}
 abs,err:=ProjectPath([]string{root})
 if err!=nil||abs!=root{t.Fatalf("absolute path: %q %v",abs,err)}
 if _,err:=ProjectPath([]string{root,root});err==nil{t.Fatal("accepted multiple projects")}
 if _,err:=ProjectPath([]string{filepath.Join(root,"missing")});err==nil{t.Fatal("accepted non-existent project")}
 file:=filepath.Join(root,"config")
 if err:=os.WriteFile(file,[]byte("test"),0600);err!=nil{t.Fatal(err)}
 if _,err:=ProjectPath([]string{file});err==nil{t.Fatal("accepted file as project")}
 symlink:=filepath.Join(t.TempDir(),"alias")
 if err:=os.Symlink(root,symlink);err==nil {
  if _,err:=ProjectPath([]string{symlink});err==nil{t.Fatal("accepted symlink")}
 }
 if _,err:=ProjectPath([]string{"--help"});!errors.Is(err,flag.ErrHelp){t.Fatalf("unexpected help response: %v",err)}
}
func TestRelativePath(t *testing.T){
 root:=t.TempDir()
 old,err:=os.Getwd();if err!=nil{t.Fatal(err)}
 if err:=os.Chdir(root);err!=nil{t.Fatal(err)}
 defer os.Chdir(old)
 got,err:=ProjectPath([]string{"."})
 if err!=nil||got!=root{t.Fatalf("dot path: %q %v",got,err)}
}
