# PKU-PPoV

- #### 介绍
     PPoV_Blockchain是基于PPoV（Parallel Proof of Vote）算法的区块链平台，以GO语言开发而成，支持多标识网络体系MIN的标识管理和日志记录等功能

- #### 运行
    -  快速启动
    
    ```shell script
    ./pku_ppov
    ```
     -  一般形式

    ```shell script
    ./pku_ppov -f <config_file>
    ```
 
    >如果无法运行请先检查MongoDB数据库服务是否开启，如未开启可使用以下命令开启:                                                                                                                                       
    ```shell script
    sudo service mongod start
    ```   
    >如果MongoDB尚未安装，请参考如下命令进行安装:                                                                                                                                       
    ```shell script
   sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 \
   --recv 0C49F3730359A14518585931BC711F9BA15703C6
  
   echo "deb [ arch=amd64,arm64 ] http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.4 multiverse" \ | 
   sudo tee /etc/apt/sources.list.d/mongodb-org-3.4.list
  
   sudo apt-get update && sudo apt-get install -y mongodb-org
   ``` 
- #### 开发环境安装
    -  使用到的GO插件
    
    ```shell script
    # 执行以下命令建议开启VPN
    go get -u -t github.com/tinylib/msgp
    go get github.com/robfig/config
    go get github.com/larspensjo/config
    go get -u github.com/tjfoc/gmsm
    go get gopkg.in/mgo.v2
    ```
      
- #### 使用说明
     默认（快速启动）情况下，在本地运行四个区块链节点，IP地址为127.0.0.1，区块链通信端口`5010, 5011, 5012, 5013`,对外提供服务端口`8010, 8011, 8012, 8013`

- #### 模块划分
   ```textmate
     AccountManager  |   账号管理模块
          bin        |   二进制可执行文件
      ConfigHelper   |   配置文件生成模块
           lib       |   使用的第三方库
         Message     |   消息管理模块
        MetaData     |   元数据格式定义
         MongoDB     |   数据库模块
       ndnKeyChain   |   MIN证书管理模块
         Network     |   网络模块
          Node       |   程序核心模块
          tools      |   集群管理模块
          utils      |   常用工具模块
    ```

- #### 配置文件说明
  ```textmate
  [network]
  "WorkerList" :            记账节点IP与端口
  "WorkerCandidateList" :   候选记账节点IP与端口
  "VoterList" :             投票节点IP与端口
  "SingleServerNodeNum" :   本机运行的节点数
  "IP" :                    本机IP地址
  "Port" :                  本机端口
  "ServicePort" :           对外提供服务端口,RPC端口
  "ManagementServerIP" :    后台服务器IP
  "ManagementServerPort" :  后台服务器端口
  "Hostname" :              本机节点名
  
  [Consensus] 
  "PubkeyList" :            本机运行节点的公钥
  "PrikeyList" :            本机运行节点的私钥
  "MyPubkey" :              保留字段
  "MyPrikey" :              保留字段
  "GenesisDutyWorker" :     生成创世区块的节点编号
  "WorkerNum" :             记账节点数
  "VotedNum" :              投票节点数
  "BlockGroupPerCycle" :    记账节点轮换周期
  "Tcut" :                  超时时间
  "GenerateBlockPeriod" :   产生区块周期
  "TxPoolSize" :            交易池大小
  "ByzantineNode" :         是否为拜占庭节点
  "DropDatabase" :          清空数据库
  
  ```
