'use strict';
/**
 * Created by Administrator on 2016/9/7.
 */

//开启路由器
app.controller("sbglGzflContr", ['$scope', 'httpService', '$modal', 'toaster', function($scope, httpService, $modal, toaster) {
	console.log("设备管理-故障分类");
	$scope.sbxh_data = {};
	$scope.sbxh_data_tree = [];

	//   故障分类变量
	$scope.gzflRightData = {
		//   设备型号选中项
		"item": {},
		//   故障分类列表数组
		"list": []
	};

	//    分页
	$scope.backPage = {
		//   超始页
		PageIndex: 1,
		//   每页显示
		PageSize: 15,
		//   页码显示条数
		"pageNumber": 5
	}

	$scope.form = {
		"KeyWord": "",
		"ModelId": ""
	}
	var treeIsOpen = [];

	//获取故障分类列表
	var getDeviceModelFaultTypeList = function(ModelId) {
		//	接口数据
		var url = config.HttpUrl + "/device/getDeviceModelFaultTypeList";
		//存接口传进来的数据
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: $scope.backPage.PageIndex,
				PageSize: $scope.backPage.PageSize
			},
			Para: {
				KeyWord: $scope.form.KeyWord,
				ModelId: ModelId
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(Data) {
			console.log("获取故障分类列表", Data);
			if(Data.Rcode == "1000") {

				$scope.gzflRightData.list = Data.Result.Data;
				//	  分页
				$scope.backPage = pageFn(Data.Result.Page, $scope.backPage.pageNumber);

			} else {
				$scope.gzflRightData.list = [];
				console.log(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//获取设备型号列表
	var getDeviceModelList = function() {
			//   是展开时加入展开数组
			treeIsOpen = [];
			dg_tree($scope.sbxh_data_tree, function(item) {
				if(item.expanded == true && item.level == $scope.sbxhRightData.item.level) {
					treeIsOpen.push({
						"level": item.level,
						"Id": item.Id,
						"expanded": true
					});
				}
			});
			//dg_tree($scope.sbxh_data_tree,function(item){
			//	if(item.expanded == true){
			//		switch(item.level){
			//			case 1:
			//				treeIsOpen.push({"level":item.level,"Id":item.Id,"expanded":true});
			//				break;
			//			case 2:
			//				treeIsOpen.push({"level":item.level,"Id":item.Id,"expanded":true});
			//				break;
			//		}
			//	}
			//});
			//	接口数据
			var url = config.HttpUrl + "/device/getDeviceModelList";
			//存接口传进来的数据
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Page: {
					PageIndex: -1,

				},
				Para: {
					KeyWord: $scope.form.KeyWord,
					ModelId: $scope.form.ModelId
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(Data) {
				console.log("设备型号查询", Data);
				if(Data.Rcode == "1000") {
					$scope.sbxh_data = Data.Result.Data;
					$scope.sbxh_data_tree = outTree($scope.sbxh_data);
					$scope.sbxh_data_tree[0].expanded = true;
					//   加入展开
					for(var a = 0; a < treeIsOpen.length; a++) {
						dg_tree($scope.sbxh_data_tree, function(item) {
							if(item.level == treeIsOpen[a].level && item.Id == treeIsOpen[a].Id) {
								item.expanded = true;
							}
						});
					}
				} else {
					console.log(data.Reason);
				}
			}, function(reason) {

			}, function(update) {});
		}
		//  递归树
	var dg_tree = function(tree, fn) {
			for(var i = 0; i < tree.length; i++) {
				fn(tree[i]);
				//			console.log(tree[i],tree[i].children)
				if(tree[i].children.length > 0) {
					dg_tree(tree[i].children, fn);
				}
			}
		}
		/*    生成节点型号树      */
		//   []
	var outTree = function(det) {
			var tree = {};
			console.log(tree);
			for(var a in det) {
				var item = det[a];
				//console.log("item",item);
				item.label = item.Name;
				item.level = item.Type;

				if(!tree[item.Id]) {
					tree[item.Id] = {};
				}

				tree[item.Id] = $.extend({}, tree[item.Id], item);
				if(!("children" in tree[item.Id])) tree[item.Id].children = [];

				if(tree[item.PId]) {
					tree[item.PId].children.push(tree[item.Id]);
				} else {
					tree[item.PId] = {
						children: [tree[item.Id]]
					};
				}
			}
			return [{
				label: '全部设备分类型号',
				children: tree[""].children,
				level: 0
			}];
		}
		//左边树 点击
	$scope.gzfl_tree_handler = function(branch) {
		console.log("ss", branch);
		//$scope.page.index = 1;
		//$scope.page.oneSize = 15;
		$scope.gzflRightData.item = branch;
		//$scope.gzflRightData.list=[];
		//		for(var c in $scope.gzfl_data){
		//			if($scope.gzflRightData.item.Id == $scope.gzfl_data[c].ModelId){
		//				$scope.gzflRightData.list.push($scope.gzfl_data[c]);
		//			}
		//		}
		getDeviceModelFaultTypeList($scope.gzflRightData.item.Id);

	};
	/*  -------------------- 分页、页码  -----------------------  */
	//$scope.backPage = {};
	/*----------------
	 //    分页对象添加页码
	 //    return  obj  分页对象
	 //    pagedata:obj  分页对象
	 //    maxpagenumber:int  显示页码数默认5个页码
	 ------------------*/
	var pageFn = function(pagedata, maxpagenumber) {
		if(pagedata.length < 1) return null;
		//   缺省时分5页
		Number(maxpagenumber) > 0 ? maxpagenumber = Number(maxpagenumber) : maxpagenumber = 5;
		var nub = [];
		var mid = Math.ceil(maxpagenumber / 2);
		if(pagedata.PageCount > maxpagenumber) {
			//  起始页
			var Snumber = 1;
			if((pagedata.PageIndex - mid) < 1) {
				Snumber = 1
			} else if((pagedata.PageIndex + mid) > pagedata.PageCount) {
				Snumber = pagedata.PageCount - maxpagenumber + 1;
			} else {
				Snumber = pagedata.PageIndex - (mid - 1)
			}
			for(var i = 0; i < maxpagenumber; i++) {
				nub.push(Snumber + i);
			}
		} else {
			for(var i = 0; i < pagedata.PageCount; i++) {
				nub.push(i + 1);
			}
		}
		pagedata.Number = nub;
		return pagedata;
	}

	//  翻页
	$scope.pageClick = function(pageindex) {
			if(!(Number(pageindex) > 0)) return false;
			if(pageindex > 0 && pageindex <= $scope.backPage.PageCount) {
				$scope.backPage.PageIndex = pageindex;
				getDeviceModelFaultTypeList($scope.gzflRightData.item.Id);
			}
		}
		/*  -------------------- 分页、页码  -----------------------  */
		//////////////////////////////////////////////////////////////////////
		//  删除列表
	var deleteDeviceModelFaultType = function(Id) {
			var url = config.HttpUrl + "/device/deleteDeviceModelFaultType";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					Id: Id
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("删除成功", data);
				if(data.Rcode == "1000") {
					getDeviceModelFaultTypeList($scope.gzflRightData.item.Id);
					toaster.pop('success', '删除成功！');
				} else {
					toaster.pop('warning', data.Reason);
				}
			});
		}
		//  添加按钮功能
	$scope.openModalAddFl = function(item, str) {
		if(!str) str = "";
		if(!item) item = {};
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/sbpz/modal_gzfl.html',
			controller: 'modalSbglSbpzGzflContr',
			windowClass: 'm-modal-sbgl-sbpz',
			resolve: {
				items: function() {
					return {
						"operate": str,
						"item": item
					};
				}
			}
		});
		modalInstance.result.then(function(bul) {
			console.log("bul", bul);
			if(bul) {
				getDeviceModelFaultTypeList($scope.gzflRightData.item.Id);
			}
		});
	}

	//  删除前提示
	var deleteBefore = function(Id, delFn) {
		var url = config.HttpUrl + "/device/onDeletingDeviceModelFaultType";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				Id: Id
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("删除前提示", data);
			if(data.Rcode == "1000") {
				delFn(data.Result);
			} else {
				return false;
			}
		});
	}

	$scope.deleteItem = function(item) {

		var modalInstance = $modal.open({
			templateUrl: 'modal/modal_alert_all.html',
			controller: 'modalAlert2Conter',
			resolve: {
				items: function() {
					return {
						"type": 'warning',
						"msg": '确定删除当前故障分类吗？'
					};
				}
			}
		});
		modalInstance.result.then(function(bul) {
			console.log("bul", bul);
			if(bul) {
				deleteBefore(item.Id, function(rBul) {
					//   是否有关联内容
					if(rBul) {
						var modalInstance2 = $modal.open({
							templateUrl: 'modal/modal_alert_all.html',
							controller: 'modalAlert2Conter',
							resolve: {
								items: function() {
									return {
										"type": 'warning',
										"msg": '警告：此操作将影响所有与当前故障分类关联的设备故障数据完整性！'
									};
								}
							}
						});
						modalInstance2.result.then(function(bul) {
							console.log("bul", bul);
							if(bul) {
								//    删除
								deleteDeviceModelFaultType(item.Id);
							}
						});
					} else {
						//  删除
						deleteDeviceModelFaultType(item.Id);
					}
				});
			}
		});

	}

	$scope.run = function() {
		//getDeviceModelFaultTypeList();
		getDeviceModelList();
	}
	$scope.run();
}]);

app.controller("modalSbglSbpzGzflContr", ['$scope', 'httpService', '$modal', '$modalInstance', 'items', 'formValidate', 'toaster', function($scope, httpService, $modal, $modalInstance, items, formValidate, toaster) {
	console.log("设备管理-设备配置-故障分类-添加故障分类弹窗");

	$scope.items = items;
	console.log("弹出items", items);

	$scope.form = {
			"Id": "",
			"Name": "",
			"ModelId": "",
			"ModelName": ""
		}
		//   生成字符串GUID
	function getGUIDs() {
		var GUID = "";
		for(var i = 1; i <= 32; i++) {
			var n = Math.floor(Math.random() * 16.0).toString(16);
			GUID += n;
			if((i == 8) || (i == 12) || (i == 16) || (i == 20))
				GUID += "";
		}
		GUID += "";
		return GUID;
	}
	//故障分类--添加
	var saveDeviceModelFaultTypeAdd = function() {
			if(!(formValidate($scope.form.Name).minLength(0).outMsg(2522).isOk)) return false;
			var url = config.HttpUrl + "/device/saveDeviceModelFaultType";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					"Id": $scope.form.Id,
					"Name": $scope.form.Name,
					"ModelName": $scope.form.ModelName,
					"ModelId": $scope.form.ModelId
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("控制命令-添加", data);
				if(data.Rcode == "1000") {
					toaster.pop('success', "添加成功！");
					$modalInstance.close(true);
				} else {
					toaster.pop('warning', data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//故障分类--修改
	var saveDeviceModelFaultType = function(Id) {
			if(!(formValidate($scope.form.Name).minLength(0).outMsg(2522).isOk)) return false;
			if(!Id) return false;
			var url = config.HttpUrl + "/device/saveDeviceModelFaultType";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					"Id": $scope.form.Id,
					"Name": $scope.form.Name,
					"ModelId": $scope.form.ModelId
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("控制命令修改", data);
				if(data.Rcode == "1000") {
					toaster.pop('success', "修改成功！");
					$modalInstance.close(true);
				} else {
					toaster.pop('warning', data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//保存按钮
	$scope.ok = function() {
			if($scope.items.operate == "edit") {
				saveDeviceModelFaultType($scope.items.item.Id);
			} else {
				saveDeviceModelFaultTypeAdd();
			}
		}
		//	取消按钮
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	}
	$scope.run = function() {
		if(items.operate == "edit" || items.operate == "see") {
			$scope.form = $.extend({}, $scope.form, $scope.items.item);
			$scope.form.ModelName = $scope.items.item.ModelName;
		} else {
			$scope.form.Id = getGUIDs();
			$scope.form.ModelId = $scope.items.item.Id;
			$scope.form.ModelName = $scope.items.item.Name;
		}
	}
	$scope.run();
}]);