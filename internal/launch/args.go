package launch

import (
 "errors"
 "flag"
 "fmt"
 "io"
 "os"
 "path/filepath"
)

// ProjectPath resolves the optional positional project path.
// An empty result means AgentForge starts without a selected project.
func ProjectPath(args []string) (string,error) {
 flags:=flag.NewFlagSet("agentforge",flag.ContinueOnError)
 flags.SetOutput(io.Discard)
 help:=flags.Bool("help",false,"show usage")
 if err:=flags.Parse(args);err!=nil {return "",fmt.Errorf("invalid arguments: %w",err)}
 if *help {return "",flag.ErrHelp}
 if flags.NArg()>1{return "",errors.New("expected at most one project directory")}
 if flags.NArg()==0{return "",nil}
 selected:=flags.Arg(0)
 if selected=="" {return "",errors.New("project directory cannot be empty")}
 abs,err:=filepath.Abs(selected);if err!=nil{return "",err}
 info,err:=os.Lstat(abs);if err!=nil{return "",fmt.Errorf("cannot open project %q: %w",selected,err)}
 if info.Mode()&os.ModeSymlink!=0||!info.IsDir(){return "",errors.New("project must be a real directory, not a symlink or file")}
 return abs,nil
}
