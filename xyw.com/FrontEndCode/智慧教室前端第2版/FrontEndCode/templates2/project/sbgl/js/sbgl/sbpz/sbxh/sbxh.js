// controller  开启路由器
app.controller("sbglSbxhContr", ['$scope', 'httpService', '$modal', 'toaster', function($scope, httpService, $modal, toaster) {
	console.log("设备配置-设备型号");
	$scope.sbxh_data = {};
	//   设备分类树不含型号
	$scope.sbxh_data_tree = [];
	//   设备分类型号树
	$scope.sbxhfl_data_tree_noxh = [];
	$scope.sbxhRightData = {
		"item": {},
		"list": []
	};
	$scope.page = {
		//   超始页
		"index": 1,
		//   每页显示
		"oneSize": 15,
		//   页码显示条数
		"pageNumber": 5
	}
	$scope.form = {
		"KeyWord": "",
		"ModelId": ""
	}

	//   树 已展开 保存数组
	var treeIsOpen = [];

	var getDeviceModelList = function(PageIndex, PageSize) {
		Number(PageIndex) > 0 ? PageIndex = Number(PageIndex) : PageIndex = 1;
		Number(PageSize) > 0 ? PageSize = Number(PageSize) : PageSize = 15;
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
		//接口路径
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
				PageIndex: -1
					//PageSize:PageSize
			},
			Para: {
				KeyWord: $scope.form.KeyWord,
				ModelId: $scope.form.ModelId
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(Data) {
			console.log("getall", Data);
			if(Data.Rcode == "1000") {
				$scope.sbxh_data = Data.Result.Data;
				$scope.sbxh_data_tree = outTree2($scope.sbxh_data);
				console.log("树", $scope.sbxh_data_tree);
				$scope.sbxh_data_tree[0].expanded = true;
				//   分类型号树
				$scope.sbxhfl_data_tree = outTree($scope.sbxh_data);
				//	分页
				$scope.backPage = pageFn(Data.Result.Page, $scope.page.pageNumber);
				//	刷新左边列表
				dg_tree($scope.sbxh_data_tree, function(item) {
					if($scope.sbxhRightData.item.Id == item.Id) {
						item.selected = true;
					}
				});
				//	刷新右边列表
				dg_tree($scope.sbxhfl_data_tree, function(item) {
					if($scope.sbxhRightData.item.Id == item.Id) {
						$scope.sbxhRightData.list = item.children;
					}
				});
				//   加入展开
				for(var a = 0; a < treeIsOpen.length; a++) {
					dg_tree($scope.sbxh_data_tree, function(item) {
						if(item.level == treeIsOpen[a].level && item.Id == treeIsOpen[a].Id) {
							item.expanded = true;
						}
					});
				}
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
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
			label: '全部设备分类',
			children: tree[""].children,
			level: 0
		}];
	}

	/*    生成节点型号树  不含型号     */
	//   []
	var outTree2 = function(det) {
		var tree = {};
		console.log(tree);
		for(var a in det) {
			var item = det[a];
			//  排除型号
			if(item.Type == '2') continue;
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
			label: '全部设备分类',
			children: tree[""].children,
			level: 0
		}];
	}

	/*  -------------------- 分页、页码  -----------------------  */
	$scope.backPage = {};
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
				getDeviceModelList(pageindex, $scope.page.oneSize);
			}
		}
		/*  -------------------- 分页、页码  -----------------------  */
		//////////////////////////////////////////////////////////////////////
		//  删除列表
	var deleteDeviceModel = function(Id) {
			var url = config.HttpUrl + "/device/deleteDeviceModel";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					ModelId: Id
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("删除成功", data);
				if(data.Rcode == "1000") {
					getDeviceModelList();
					toaster.pop('success', '删除成功！');
				} else {
					toaster.pop('warning', data.Reason);
				}
			});
		}
		//左边树 点击
	$scope.sbxh_tree_handler = function(branch) {
		//console.log("单机",$scope.sbxh_data_tree);
		//console.log("单机",branch);
		$scope.page.index = 1;
		$scope.page.oneSize = 15;
		$scope.sbxhRightData.item = branch;
		//
		dg_tree($scope.sbxh_data_tree, function(item) {
			if(branch.Id == item.Id) {
				item.selected = true;
			} else {
				item.selected = false;
			}
		});
		//   查找分类型号树中 children 属性数据
		dg_tree($scope.sbxhfl_data_tree, function(item) {
			if(branch.Id == item.Id) {
				$scope.sbxhRightData.list = item.children;
			}
		});
	};

	//  添加按钮功能
	$scope.openModalAddXh = function(item, str) {
			if(!str) str = "";
			if(!item) item = {};
			var modalInstance = $modal.open({
				templateUrl: '../project/sbgl/html/sbgl/sbpz/modal_sbxh.html',
				controller: 'modalSbglSbpzSbxhContr',
				windowClass: 'm-modal-sbgl-sbpz',
				resolve: {
					items: function() {
						item.ModalName = $scope.sbxhRightData.item.Name;
						//item.ModalId = $scope.sbxhRightData.item.Id;
						item.ModalPId = $scope.sbxhRightData.item.Id;
						item.Modallevel = $scope.sbxhRightData.item.level;
						return {
							"operate": str,
							"item": item,
							"enumDeviceModel": $scope.enumDeviceModel
						};
					}
				}
			});
			modalInstance.result.then(function(bul) {
				console.log("bul", bul);
				if(bul) {
					getDeviceModelList();
				}
			});
		}
	
	
	
	//  删除前提示
	var deleteBefore = function(Id,delFn) {
		var url = config.HttpUrl + "/device/onDeletingDeviceModel";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				ModelId: Id
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
	
	
	
	
	
	
	//  删除弹窗功能
	$scope.deleteItem = function(item) {
		
		//   设备型号提示
		var msg_xh = ['确定删除当前设备型号吗？','警告：此操作将影响所有与当前设备型号关联的设备控制、图片资源显示、故障现象词条、故障分类！'];
		//   设备分类提示
		var msg_fl = ['确定删除当前设备分类吗？','警告：此操作将删除其下所有设备型号！'];
		//
		var msg_txt = [];
		if(item.Type == 1){
			//   分类
			msg_txt = msg_fl;
		}else if(item.Type ==2){
			//   型号
			msg_txt = msg_xh;
		}
		
		var modalInstance = $modal.open({
			templateUrl: 'modal/modal_alert_all.html',
			controller: 'modalAlert2Conter',
			resolve: {
				items: function() {
					return {
						"type": 'warning',
						"msg": msg_txt[0]
					};
				}
			}
		});
		modalInstance.result.then(function(bul) {
			console.log("bul", bul);
			if(bul) {
				deleteBefore(item.Id, function(rBul) {
					//   是否有关联内容
					if(rBul){
						var modalInstance2 = $modal.open({
							templateUrl: 'modal/modal_alert_all.html',
							controller: 'modalAlert2Conter',
							resolve: {
								items: function() {
									return {
										"type": 'warning',
										"msg": msg_txt[1]
									};
								}
							}
						});
						modalInstance2.result.then(function(bul) {
							console.log("bul", bul);
							if(bul) {
								//    删除
								deleteDeviceModel(item.Id);
							}
						});
					}else{
						//  删除
						deleteDeviceModel(item.Id);
					}
				});
			}
		}); 
		

	}
	
	//
	$scope.run = function() {
		getDeviceModelList();
	}
	$scope.run();

}]);

//设备管理-设备配置-设备型号-添加设备型号弹窗
app.controller("modalSbglSbpzSbxhContr", ['$scope', 'httpService', '$modal', '$modalInstance', 'items', 'formValidate', 'toaster', function($scope, httpService, $modal, $modalInstance, items, formValidate, toaster) {
	console.log("设备管理-设备配置-设备型号-添加设备型号弹窗 ");
	$scope.items = items;
	console.log("弹窗", $scope.items);
	$scope.form = {
			"Id": "",
			"PId": "",
			"ModalName": "",
			"ModalId": "",
			"ModalPId": "",
			"Modallevel": "",
			"Name": "",
			"Description": "",
			"Type": 1,
			"TypeName": "",
			"PageFileName": "",
			"ImgFileName": "",
			"ImgFileName2": "",
			//    是否按时间预警[enum=0:不预警/1:预警]
			"IsAlert": "1",
			//    设置超出的时间值(秒)
			"MaxUseTime": 0,
			//  选择类型状态
			"TypeStates": [{
				value: "1",
				Name: "分类"
			}, {
				value: "2",
				Name: "型号"
			}],
			//  类型状态默认参数
			"TypeState": "",
			//  是否按时间预警
			"IsAlerts": [{
				value: "0",
				Name: "否"
			}, {
				value: "1",
				Name: "是"
			}],
			//  是否按时间预警默认参数
			"IsAler": "",
			//  资源模板选择数组
			"enumDeviceArrays": []
		}
		//  把类型选中状态的值传给 Type
	$scope.changeTypeState = function(item) {
			$scope.form.Type = item.value;
		}
		//  把按时间预警的值传给 IsAlert
	$scope.changeIsAlert = function(item) {
			$scope.form.IsAlert = item.value;
		}
		//  资源模板对象转数组
	$scope.transform = function(obj) {
			$scope.enumDeviceArrays = [];
			for(item in obj) {
				$scope.enumDeviceArrays.push(obj[item]);
			}
			return $scope.enumDeviceArrays;
		}
		//   资源模板选择
	$scope.changeEnumDeviceModel = function(str) {
		angular.forEach($scope.items.enumDeviceModel, function(val) {
			if(val.title == str.title) {
				$scope.form.PageFileName = val.PageFileName;
				$scope.form.ImgFileName = val.ImgFileName;
				$scope.form.ImgFileName2 = val.ImgFileName2;
			}
		});
	}

	//打开窗口  设备型号
	$scope.modalOpenDevice = function() {
			var modalInstance = $modal.open({
				templateUrl: '../html/modal/modal_device.html',
				controller: 'modalGetDeviceCtrl',
				resolve: {
					items: function() {
						return $scope.items;
					}
				}
			});

			modalInstance.result.then(function(deviceItem) {
				console.log("ss", deviceItem)
				if(!deviceItem) {
					$scope.form.ModalName = "";
					$scope.form.ModalPId = "";
				} else {
					$scope.form.ModalName = deviceItem.Name;
					$scope.form.ModalPId = deviceItem.Id;
				}
			}, function() {
				//$log.info('Modal dismissed at: ' + new Date());
			});
		}
		//设备型号--添加
	var saveDeviceModelAdd = function() {
			if(!(formValidate($scope.form.Id).minLength(0).outMsg(2500).isOk)) return false;
			if(!(formValidate($scope.form.Name).minLength(0).outMsg(2501).isOk)) return false;
			if(!(formValidate($scope.form.Type).isNumber().outMsg(2502).isOk)) return false;

			//   小时转秒
			$scope.form.MaxUseTime = Number($scope.form.MaxUseTime) * 3600;

			var url = config.HttpUrl + "/device/saveDeviceModel";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					"Id": $scope.form.Id,
					"PId": $scope.form.ModalPId,
					"Name": $scope.form.Name,
					"Description": $scope.form.Description,
					"Type": Number($scope.form.Type),
					"PageFileName": $scope.form.PageFileName,
					//"TypeName":$scope.form.TypeName,
					"ImgFileName": $scope.form.ImgFileName,
					"ImgFileName2": $scope.form.ImgFileName2,
					"IsAlert": $scope.form.IsAlert,
					"MaxUseTime": $scope.form.MaxUseTime
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("设备型号--添加", data);
				if(data.Rcode == "1000") {
					toaster.pop('success', '添加成功！');
					$modalInstance.close(true);
				} else {
					toaster.pop('warning', data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//设备型号--修改
	var saveDeviceModel = function(Id) {

			if(!(formValidate($scope.form.Id).minLength(0).outMsg(2500).isOk)) return false;
			if(!(formValidate($scope.form.Name).minLength(0).outMsg(2501).isOk)) return false;
			if(!(formValidate($scope.form.Type.toString()).minLength(0).outMsg(2502).isOk)) return false;

			if(!Id) return false;
			var url = config.HttpUrl + "/device/saveDeviceModel";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					"Id": $scope.items.item.Id,
					"PId": $scope.form.ModalPId,
					"Name": $scope.form.Name,
					"Description": $scope.form.Description,
					"Type": Number($scope.form.Type),
					"PageFileName": $scope.form.PageFileName,
					//"TypeName":$scope.form.TypeName,
					"ImgFileName": $scope.form.ImgFileName,
					"ImgFileName2": $scope.form.ImgFileName2,
					"IsAlert": $scope.form.IsAlert,
					"MaxUseTime": $scope.form.MaxUseTime

				}
			}
			console.log("修改", $scope.form.Type);
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("控制命令修改", data);
				if(data.Rcode == "1000") {
					toaster.pop('success', '修改成功！');
					$modalInstance.close(true);
				} else {
					toaster.pop('warning', data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//保存按钮
	$scope.ok = function() {
			if($scope.items.operate == "edit") {
				saveDeviceModel($scope.items.item.Id);
			} else {
				saveDeviceModelAdd();
			}
		}
		//   生成字符串GUID
	function getGUIDs() {
		var GUID = "";
		for(var i = 1; i <= 32; i++) {
			var n = Math.floor(Math.random() * 16.0).toString(16);
			GUID += n;
			if((i == 8) || (i == 12) || (i == 16) || (i == 20))
			//GUID += "-";
				GUID += "";
		}
		GUID += "";
		return GUID;
	}

	//	取消按钮
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	}

	$scope.run = function() {
		$scope.form.enumDeviceArrays = $scope.transform($scope.items.enumDeviceModel);
		if(items.operate == "edit" || items.operate == "see") {
			$scope.form = $.extend({}, $scope.form, $scope.items.item);
			$scope.form.ModalName = $scope.items.item.ModalName;
			console.log("查看修改", $scope.items);
			for(var item in $scope.form.TypeStates) {
				if($scope.form.TypeStates[item].value == $scope.form.Type) {
					$scope.form.TypeState = $scope.form.TypeStates[item];
				}
			}
			for(var item in $scope.form.IsAlerts) {
				if($scope.form.IsAlerts[item].value == $scope.form.IsAlert) {
					$scope.form.IsAler = $scope.form.IsAlerts[item];
				}
			}
			for(var item in $scope.form.enumDeviceArrays) {
				if($scope.form.enumDeviceArrays[item].ImgFileName == $scope.form.ImgFileName) {
					$scope.form.enumDeviceArray = $scope.form.enumDeviceArrays[item];
				}
			}
		} else {
			$scope.form.Id = getGUIDs();
			$scope.form.ModalPId = $scope.items.item.ModalPId;
			$scope.form.ModalName = $scope.items.item.ModalName;
			$scope.form.Modallevel = $scope.items.item.Modallevel;
			$scope.form.TypeState = $scope.form.TypeStates[0];
			$scope.form.IsAler = $scope.form.IsAlerts[1];
			if($scope.items.item.Modallevel == undefined || $scope.items.item.Modallevel == 0) {
				$scope.form.TypeStates.pop();
			}
		}
	}
	$scope.run();
}]);
/**
 * Created by liuqi on 2016/12/5.
 */