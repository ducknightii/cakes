package main

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"math/rand"
	"time"
)

// [[0.6046603 0.9405091] [0.84640664 0.7586615]]
var collName = "test"
var partitionName = "part5"
var milClient client.Client
var pickedVec [][]float32

func main() {
	var err error
	//...other snippet ...
	milClient, err = client.NewGrpcClient(context.Background(), "172.16.17.18:19530")
	defer milClient.Close()
	if err != nil {
		// handle error
	}

	// check exist
	collExist, err := milClient.HasCollection(context.Background(), collName)
	fmt.Printf("collection: %s exist:%v err:%v\n", collName, collExist, err)
	/*if collExist {
		if err = dropCollection(); err != nil {
			return
		}
	}*/

	// create collection
	// notes: if exist return err: collection test already exist
	/*err = createCollection()
	if err != nil {
		return
	}*/

	// desc collection
	if err = descCollection(); err != nil {
		return
	}

	// has partition
	if exist, _ := hasPartition(); !exist {
		// create partition
		if err = createPartition(); err != nil {
			return
		}
	}

	// notes: 根据 插入 删除 比对的现象猜测，不采用自增主键 而是外部传入的方式(uuid) 对于milvus并不是真的唯一主键，而是只是个用于删除返回查询结果的字段
	// 实验中  多次使用0-1999键值插入数据，然后再删除1，1111两个数据，然后以最新一次的 0 1111向量进行查询，第一次查不出来1111，后面却始终可以查出1111，这是采用标记删除，发现一个索引是1111的已经删除就直接返回了吧，实际只有1111的数据存在呢
	// insert data
	// notes: return success even data exist
	/*if err = insertData(); err != nil {
		return
	}*/

	// del data
	// notes: return success even data not exist
	if err = delData(); err != nil {
		return
	}
	fmt.Printf("insert vec: %v\n", pickedVec)
	pickedVec = [][]float32{
		[]float32{0.6046603, 0.9405091},
		[]float32{0.84640664, 0.7586615},
	}

	// load collection
	// notes: 没有数据时 异常退出了
	/*if err = loadCollection(); err != nil {
		return
	}*/

	// search
	if err = search(); err != nil {
		return
	}
	time.Sleep(time.Second * 15)
	if err = search(); err != nil {
		return
	}
}

// create collection
func createCollection() error {
	schema := entity.Schema{
		CollectionName: collName,
		Description:    "测试",
		AutoID:         false,
		Fields: []*entity.Field{
			&entity.Field{
				ID:          0,
				Name:        "uuid",
				PrimaryKey:  true,
				AutoID:      false,
				Description: "unique ",
				DataType:    entity.FieldTypeInt64, // 只有int64 可以做主键
			},
			/*&entity.Field{       // string field not support
				ID:          0,
				Name:        "user_id",
				PrimaryKey:  false,
				AutoID:      false,
				Description: "user id",
				DataType:    entity.FieldTypeString,
			},*/
			&entity.Field{
				ID:          0,
				Name:        "feature",
				PrimaryKey:  false,
				AutoID:      false,
				Description: "feature",
				DataType:    entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "2",
				},
			},
		},
	}
	err := milClient.CreateCollection(context.Background(), &schema, 2)
	if err != nil {
		fmt.Printf("create collection: %s err: %s\n", collName, err)
	} else {
		fmt.Printf("create collection: %s success\n", collName)
	}
	return err
}

// describe collection
func descCollection() error {
	collDesc, err := milClient.DescribeCollection(context.Background(), collName)
	if err != nil {
		fmt.Printf("describe collection: %s err: %s\n", collName, err)
	} else {
		fmt.Printf("collection: %s info: %+v\n", collName, collDesc)
	}

	return err
}

// del collection
func dropCollection() error {
	err := milClient.DropCollection(context.Background(), collName)
	if err != nil {
		fmt.Printf("collection: %s drop err: %s\n", collName, err)
	} else {
		fmt.Printf("collection: %s drop success\n", collName)
	}
	return err
}

// load collection
func loadCollection() error {
	err := milClient.LoadCollection(context.Background(), collName, false)
	if err != nil {
		fmt.Printf("load collection: %s err: %s\n", collName, err)
	} else {
		fmt.Printf("load collection: %s success\n", collName)
	}

	return err
}

// create partition
func createPartition() error {
	err := milClient.CreatePartition(context.Background(), collName, partitionName)
	if err != nil {
		fmt.Printf("create partition: %s err: %s\n", partitionName, err)
	} else {
		fmt.Printf("create partition: %s success\n", partitionName)
	}

	return err
}

// has partition
func hasPartition() (exist bool, err error) {
	exist, err = milClient.HasPartition(context.Background(), collName, partitionName)
	if err != nil {
		fmt.Printf("partition: %s check err: %s\n", partitionName, err)
	} else {
		fmt.Printf("partition: %s exist\n", partitionName)
	}
	return
}

// insert data
func insertData() error {
	rand.Seed(time.Now().Unix())
	bookIDs := make([]int64, 0, 2000)
	//userIDs := make([]string, 0, 2000)  // string field not support
	bookIntros := make([][]float32, 0, 2000)
	for i := 0; i < 2000; i++ {
		bookIDs = append(bookIDs, int64(i))
		//userIDs = append(userIDs, strconv.Itoa(i+10000))
		v := make([]float32, 0, 2)
		for j := 0; j < 2; j++ {
			v = append(v, rand.Float32())
		}
		if i%1111 == 0 {
			pickedVec = append(pickedVec, v)
		}
		bookIntros = append(bookIntros, v)
	}
	idColumn := entity.NewColumnInt64("uuid", bookIDs)
	//userColumn := entity.NewColumnString("user_id", userIDs)
	introColumn := entity.NewColumnFloatVector("feature", 2, bookIntros)

	_, err := milClient.Insert(
		context.Background(), // ctx
		collName,             // CollectionName
		partitionName,        // partitionName
		idColumn,             // columnarData
		//userColumn,         // columnarData
		introColumn, // columnarData
	)
	if err != nil {
		fmt.Printf("failed to insert data: %s\n", err)
	} else {
		fmt.Printf("insert data success\n")
	}

	return err
}

// search
func search() error {
	sp, _ := entity.NewIndexFlatSearchParam(1000)
	var vecs []entity.Vector
	for _, _vec := range pickedVec {
		vecs = append(vecs, entity.FloatVector(_vec))
	}
	fmt.Printf("search vec len: %d val: %+v\n", len(vecs), vecs)
	res, err := milClient.Search(context.Background(),
		collName,
		[]string{partitionName},
		"",
		[]string{"uuid"},
		vecs,
		"feature",
		entity.L2,
		20,
		sp,
	)
	if err != nil {
		fmt.Printf("search err: %s\n", err)
	} else {
		fmt.Printf("====> search res: %d\n", len(res))
		for _, _res := range res {
			fmt.Printf("\tcount: %d\n", _res.ResultCount)
			fmt.Printf("\tids:%s len:%d\n", _res.IDs.Name(), _res.IDs.Len())
			fmt.Printf("\tfields: %d\n", len(_res.Fields))
			for _, _field := range _res.Fields {
				fmt.Printf("\t\t%s: len:%d val:%s\n", _field.Name(), _field.Len(), _field.FieldData())
			}
			fmt.Printf("\tscore:%v\n", _res.Scores)
			fmt.Println("======")
		}
	}

	return err
}

// del data
func delData() error {
	err := milClient.DeleteByPks(context.Background(), collName, partitionName, entity.NewColumnInt64("uuid", []int64{1, 1111}))
	if err != nil {
		fmt.Printf("del data err: %s\n", err)
	} else {
		fmt.Printf("del data success\n")
	}
	return err
}
