// Auto-generated by avdl-compiler v1.3.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: sserver-avdl/sserver.avdl

package searchsrv1

import (
	rpc "github.com/keybase/go-framed-msgpack-rpc"
	context "golang.org/x/net/context"
)

type DocumentID string
type FolderID string
type TlfInfo struct {
	Salts [][]byte `codec:"salts" json:"salts"`
	Size  int64    `codec:"size" json:"size"`
}

type Trapdoor struct {
	Codeword [][]byte `codec:"codeword" json:"codeword"`
}

type WriteIndexArg struct {
	TlfID       FolderID   `codec:"tlfID" json:"tlfID"`
	SecureIndex []byte     `codec:"secureIndex" json:"secureIndex"`
	DocID       DocumentID `codec:"docID" json:"docID"`
}

type RenameIndexArg struct {
	TlfID FolderID   `codec:"tlfID" json:"tlfID"`
	Orig  DocumentID `codec:"orig" json:"orig"`
	Curr  DocumentID `codec:"curr" json:"curr"`
}

type DeleteIndexArg struct {
	TlfID FolderID   `codec:"tlfID" json:"tlfID"`
	DocID DocumentID `codec:"docID" json:"docID"`
}

type GetKeyGensArg struct {
	TlfID FolderID `codec:"tlfID" json:"tlfID"`
}

type SearchWordArg struct {
	TlfID     FolderID            `codec:"tlfID" json:"tlfID"`
	Trapdoors map[string]Trapdoor `codec:"trapdoors" json:"trapdoors"`
}

type RegisterTlfIfNotExistsArg struct {
	TlfID        FolderID `codec:"tlfID" json:"tlfID"`
	LenSalt      int      `codec:"lenSalt" json:"lenSalt"`
	FpRate       float64  `codec:"fpRate" json:"fpRate"`
	NumUniqWords int64    `codec:"numUniqWords" json:"numUniqWords"`
}

type SearchServerInterface interface {
	WriteIndex(context.Context, WriteIndexArg) error
	RenameIndex(context.Context, RenameIndexArg) error
	DeleteIndex(context.Context, DeleteIndexArg) error
	GetKeyGens(context.Context, FolderID) ([]int, error)
	SearchWord(context.Context, SearchWordArg) ([]DocumentID, error)
	RegisterTlfIfNotExists(context.Context, RegisterTlfIfNotExistsArg) (TlfInfo, error)
}

func SearchServerProtocol(i SearchServerInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "searchsrv.1.searchServer",
		Methods: map[string]rpc.ServeHandlerDescription{
			"writeIndex": {
				MakeArg: func() interface{} {
					ret := make([]WriteIndexArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]WriteIndexArg)
					if !ok {
						err = rpc.NewTypeError((*[]WriteIndexArg)(nil), args)
						return
					}
					err = i.WriteIndex(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"renameIndex": {
				MakeArg: func() interface{} {
					ret := make([]RenameIndexArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]RenameIndexArg)
					if !ok {
						err = rpc.NewTypeError((*[]RenameIndexArg)(nil), args)
						return
					}
					err = i.RenameIndex(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"deleteIndex": {
				MakeArg: func() interface{} {
					ret := make([]DeleteIndexArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]DeleteIndexArg)
					if !ok {
						err = rpc.NewTypeError((*[]DeleteIndexArg)(nil), args)
						return
					}
					err = i.DeleteIndex(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"getKeyGens": {
				MakeArg: func() interface{} {
					ret := make([]GetKeyGensArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]GetKeyGensArg)
					if !ok {
						err = rpc.NewTypeError((*[]GetKeyGensArg)(nil), args)
						return
					}
					ret, err = i.GetKeyGens(ctx, (*typedArgs)[0].TlfID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"searchWord": {
				MakeArg: func() interface{} {
					ret := make([]SearchWordArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]SearchWordArg)
					if !ok {
						err = rpc.NewTypeError((*[]SearchWordArg)(nil), args)
						return
					}
					ret, err = i.SearchWord(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"registerTlfIfNotExists": {
				MakeArg: func() interface{} {
					ret := make([]RegisterTlfIfNotExistsArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]RegisterTlfIfNotExistsArg)
					if !ok {
						err = rpc.NewTypeError((*[]RegisterTlfIfNotExistsArg)(nil), args)
						return
					}
					ret, err = i.RegisterTlfIfNotExists(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type SearchServerClient struct {
	Cli rpc.GenericClient
}

func (c SearchServerClient) WriteIndex(ctx context.Context, __arg WriteIndexArg) (err error) {
	err = c.Cli.Call(ctx, "searchsrv.1.searchServer.writeIndex", []interface{}{__arg}, nil)
	return
}

func (c SearchServerClient) RenameIndex(ctx context.Context, __arg RenameIndexArg) (err error) {
	err = c.Cli.Call(ctx, "searchsrv.1.searchServer.renameIndex", []interface{}{__arg}, nil)
	return
}

func (c SearchServerClient) DeleteIndex(ctx context.Context, __arg DeleteIndexArg) (err error) {
	err = c.Cli.Call(ctx, "searchsrv.1.searchServer.deleteIndex", []interface{}{__arg}, nil)
	return
}

func (c SearchServerClient) GetKeyGens(ctx context.Context, tlfID FolderID) (res []int, err error) {
	__arg := GetKeyGensArg{TlfID: tlfID}
	err = c.Cli.Call(ctx, "searchsrv.1.searchServer.getKeyGens", []interface{}{__arg}, &res)
	return
}

func (c SearchServerClient) SearchWord(ctx context.Context, __arg SearchWordArg) (res []DocumentID, err error) {
	err = c.Cli.Call(ctx, "searchsrv.1.searchServer.searchWord", []interface{}{__arg}, &res)
	return
}

func (c SearchServerClient) RegisterTlfIfNotExists(ctx context.Context, __arg RegisterTlfIfNotExistsArg) (res TlfInfo, err error) {
	err = c.Cli.Call(ctx, "searchsrv.1.searchServer.registerTlfIfNotExists", []interface{}{__arg}, &res)
	return
}
