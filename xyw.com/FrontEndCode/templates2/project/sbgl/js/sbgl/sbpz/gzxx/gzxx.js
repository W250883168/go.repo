// controller  开启路由器
app.controller("sbglGzxxContr", ['$scope', 'httpService', '$modal', 'toaster',function ($scope, httpService, $modal,toaster) {
    console.log("设备配置-故障现象");

    //   设备列表
    $scope.FaultTypeList = [];

    //    分页
    $scope.backPage = {
        PageIndex:1,
        PageSize:10
    }

    $scope.form = {
    	//    关键词
    	"KeyWord":"",
    	//    设备型号ID
    	"ModelId":""
    }

    //    设备型号列表
    $scope.sbxh_data = {};
    //    设备型号生成后树
	$scope.sbxh_data_tree=[];
	//    右边选中后左边列表对象
	$scope.sbxhRightData = {
		"item":{},
		"list":[]
	};


    /**
     * 取故障现象列表
     */
    var getDeviceModelFaultTypeList = function(){
        var url = config.HttpUrl + "/device/getDeviceModelFaultWordList";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Page": {
				PageIndex:-1,
				//PageSize:$scope.backPage.PageSize
			},
            "Para":{
                //    查询关键字:字符串非必填
		    	"KeyWord":$scope.form.KeyWord,
		    	//    设备型号ID:字符串 非必填
		        "ModelId":$scope.form.ModelId
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("取故障现象列表",data);
            if(data.Rcode =="1000"){
                $scope.FaultTypeList = data.Result.Data;
                //    清空
				$scope.sbxhRightData.list = [];
                //   传入
				for(var a in $scope.FaultTypeList){
					if($scope.sbxhRightData.item.Id == $scope.FaultTypeList[a].ModelId){
						$scope.sbxhRightData.list.push($scope.FaultTypeList[a]);
					}
				}
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }

	/**
     * 取设备型号列表
     */
    var getDeviceModelList = function(){
    	//接口路径
    	var url=config.HttpUrl+"/device/getDeviceModelList";
    	//存接口传进来的数据
		var data = {
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: -1
				//PageSize: $scope.backPage.PageSize
			},
			Para:{
				KeyWord:$scope.form.KeyWord,
				ModelId:$scope.form.ModelId
			}
		}
    	var promise=httpService.ajaxPost(url,data);
    	promise.then(function(Data){
    		console.log("取设备型号列表",Data);
    		if(Data.Rcode=="1000"){
    			$scope.sbxh_data = Data.Result.Data;;
				$scope.sbxh_data_tree = outTree($scope.sbxh_data);
				$scope.sbxh_data_tree[0].expanded = true;
    		}else{
          toaster.pop('warning',data.Reason);
    		}
    	},function(reason){},function(update){});
    }


	/*    生成节点型号树      */
	//   []
	var outTree = function(det){
		var tree = {};
		console.log(tree);
		for(var a in det){
			var item = det[a];
			//console.log("item",item);
			item.label = item.Name;
			item.level = item.Type;

			if(!tree[item.Id]) {
				tree[item.Id] = {};
			}

			tree[item.Id] = $.extend({},tree[item.Id],item);
			if(!("children" in tree[item.Id])) tree[item.Id].children = [];

			if(tree[item.PId]){
				tree[item.PId].children.push(tree[item.Id]);
			}else{
				tree[item.PId] = {
					children: [tree[item.Id]]
				};
			}
		}
		return [{label:'全部设备分类型号',children:tree[""].children,level:0}];
	}
	//////////////////////////////////////////////////////////////////////
	//  删除列表
	var deleteDeviceModelFaultWord = function(Id){
		if(!Id)return false;
		var url = config.HttpUrl + "/device/deleteDeviceModelFaultWord";
		var data = {
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para:{
				"Id":Number(Id)
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("删除成功",data);
			if(data.Rcode =="1000"){
				getDeviceModelFaultTypeList();
        toaster.pop('success',"删除成功！");
			}else{
        toaster.pop('warning',data.Reason);
			}
		});
	}
	//左边树 点击
	$scope.sbxh_tree_handler = function(branch) {
		console.log(branch);
		$scope.backPage.PageIndex = 1;
		$scope.backPage.PageSize = 10;
		$scope.sbxhRightData.item = branch;
		//    清空
		$scope.sbxhRightData.list = [];
		//   传入
		for(var a in $scope.FaultTypeList){
			if($scope.sbxhRightData.item.Id == $scope.FaultTypeList[a].ModelId){
				$scope.sbxhRightData.list.push($scope.FaultTypeList[a]);
			}
		}
	};

	//  添加按钮功能
	$scope.openModalAddGzxx = function (str,item) {
		if(!str)str = "";
		if(!item)item = {};
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/sbpz/modal_gzxx.html',
			controller: 'modalSbglSbpzGzxxContr',
			windowClass: 'm-modal-sbgl-sbpz-gzxx',
			resolve: {
				items: function () {
					return {"operate":str,"item":item};
				}
			}
		});
		modalInstance.result.then(function(bul){
			console.log("bul",bul);
			if(bul){
				getDeviceModelFaultTypeList();
			}
		});
	}

	//   删除
	$scope.deleteItem = function(item){
		if(!window.confirm('删除故障现象后，将无法统计分析已采集此类的设备故障数据。\r\n您想继续吗'))return false;
		deleteDeviceModelFaultWord(item.Id);
	}
  $scope.deleteItem = function(item){
    var modalInstance = $modal.open({
      templateUrl: 'modal/modal_alert_all.html',
      controller: 'modalAlert2Conter',
      resolve: {
        items: function () {
          return {"type":'warning',"msg":'删除故障现象后,将无法统计分析已采集此类的设备故障数据。<br />您想继续吗？'};
        }
      }
    });
    modalInstance.result.then(function(bul){
      if(bul){
        deleteDeviceModelFaultWord(item.Id);
      }
    });
  }
	$scope.run = function(){
		getDeviceModelList();
		getDeviceModelFaultTypeList();
	}
	$scope.run();
}]);

//设备管理-设备配置-故障现象-添加故障现象弹窗
app.controller("modalSbglSbpzGzxxContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService, $modal,$modalInstance,items,formValidate,toaster) {
	console.log("设备管理-设备配置-故障现象-添加故障现象弹窗  ");

	$scope.items = items;

	//
	$scope.form = {
		//  ID
		"Id":null,
		//   故障现象名称
		"Name":"",
		//   设备型号ID
		"ModelId":"",
		//   设备型号名称
		"ModelName":""
	}


	//  添加故障现象
	var saveDeviceModelFaultWord = function(){
		if(!(formValidate($scope.form.ModelName).minLength(0).outMsg(2501).isOk))return false;
		if(!(formValidate($scope.form.Name).minLength(0).outMsg(2521).isOk))return false;

		var url = config.HttpUrl + "/device/saveDeviceModelFaultWord";
		var data = {
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para:{
				//  ID
				"Id":0,
				//   故障现象名称
				"Name":$scope.form.Name,
				//   设备型号ID
				"ModelId":$scope.form.ModelId
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("添加故障现象",data);
			if(data.Rcode =="1000"){
        toaster.pop('success',"添加成功！");
				$modalInstance.close(true);
			}else{
        toaster.pop('warning',data.Reason);
			}
		});
	}


	//  修改故障现象
	var saveDeviceModelFaultWordEdit = function(){
		if(!(formValidate($scope.form.ModelName).minLength(0).outMsg(2501).isOk))return false;
		if(!(formValidate($scope.form.Name).minLength(0).outMsg(2521).isOk))return false;
		var url = config.HttpUrl + "/device/saveDeviceModelFaultWord";
		var data = {
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para:{
				//  ID
				"Id":Number($scope.form.Id),
				//   故障现象名称
				"Name":$scope.form.Name,
				//   设备型号ID
				"ModelId":$scope.form.ModelId
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("修改故障现象",data);
      if(data.Rcode =="1000"){
        toaster.pop('success',"修改成功！");
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		});
	}


	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}

	$scope.ok = function(){
		switch($scope.items.operate){
			case "add":
				//
				saveDeviceModelFaultWord();
			break;
			case "look":
				//
			break;
			case "edit":
				//
				saveDeviceModelFaultWordEdit();
			break;
		}
	}


	$scope.run = function(){
		switch($scope.items.operate){
			case "add":
				//
				$scope.form.ModelId = $scope.items.item.Id;
				$scope.form.ModelName = $scope.items.item.Name;
			break;
			case "look":
				//
				$scope.form.Id = $scope.items.item.Id;
				$scope.form.Name = $scope.items.item.Name;
				$scope.form.ModelId = $scope.items.item.ModelId;
				$scope.form.ModelName = $scope.items.item.ModelName;
			break;
			case "edit":
				//
				$scope.form.Id = $scope.items.item.Id;
				$scope.form.Name = $scope.items.item.Name;
				$scope.form.ModelId = $scope.items.item.ModelId;
				$scope.form.ModelName = $scope.items.item.ModelName;
			break;
		}
	}
	$scope.run();

}]);
