//开启路由器
app.controller("sbglJdxhContr", ['$scope', 'httpService', '$modal', 'toaster', function($scope, httpService, $modal, toaster) {
	console.log("节点配置-节点型号");
	$scope.jdxhItems = [];
	$scope.page = {
			//   超始页
			"index": 1,
			//   每页显示
			"oneSize": 10,
			//   页码显示条数
			"pageNumber": 5
		}
		//  查询条件
	$scope.form = {
			//   搜索关键词
			"KeyWord": "",
		}
		//  查询列表
	var getNodeModelList = function(PageIndex, PageSize) {
			Number(PageIndex) > 0 ? PageIndex = Number(PageIndex) : PageIndex = 1;
			Number(PageSize) > 0 ? PageSize = Number(PageSize) : PageSize = 10;
			var url = config.HttpUrl + "/device/getNodeModelList";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Page: {
					PageIndex: PageIndex,
					PageSize: PageSize
				},
				Para: {
					KeyWord: $scope.form.KeyWord
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("接口数据 节点型号查询", data);
				if(data.Rcode == "1000") {
					$scope.jdxhItems = data.Result.Data;
					//分页
					$scope.backPage = pageFn(data.Result.Page, $scope.page.pageNumber);
				} else {
					console.log(data.Reason);
				}
			});
		}
		//  删除列表
	var deleteNodeModel = function(Id) {
		var url = config.HttpUrl + "/device/deleteNodeModel";
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
				getNodeModelList();
				toaster.pop('success', '删除成功！');
			} else {
				toaster.pop('danger', data.Reason);
			}
		});
	}

	//  删除前提示
	var deleteBefore = function(Id,delFn) {
		var url = config.HttpUrl + "/device/onDeletingNodeModel";
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

	//  添加按钮功能
	$scope.openModalAddXh = function(item, str) {
			if(!str) str = "";
			if(!item) item = {};
			var modalInstance = $modal.open({
				templateUrl: '../project/sbgl/html/sbgl/jdpz/modal_jdxh.html',
				controller: 'modalJdxhContr',
				windowClass: 'm-modal-sbgl-jdpz',
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
					getNodeModelList();
				}
			});
		}
		//  删除弹窗功能
	$scope.deleteItem = function(item) {
			//if(!window.confirm('此操作将删除其下所有控制命令，同时会解除节点与该节点型号的绑定。\r\n您想继续吗？'))return false;
			var modalInstance = $modal.open({
				templateUrl: 'modal/modal_alert_all.html',
				controller: 'modalAlert2Conter',
				resolve: {
					items: function() {
						return {
							"type": 'warning',
							"msg": '确定删除当前节点型号吗？'
						};
					}
				}
			});
			modalInstance.result.then(function(bul) {
				console.log("bul", bul);
				if(bul) {
					deleteBefore(item.Id, function(rBul) {
						if(rBul){
							var modalInstance2 = $modal.open({
								templateUrl: 'modal/modal_alert_all.html',
								controller: 'modalAlert2Conter',
								resolve: {
									items: function() {
										return {
											"type": 'warning',
											"msg": '警告：此操作将影响与当前节点型号关联的节点控制及设备控制！'
										};
									}
								}
							});
							modalInstance2.result.then(function(bul) {
								console.log("bul", bul);
								if(bul) {
									//    删除
									deleteNodeModel(item.Id);
								}
							});
						}else{
							//  删除
							deleteNodeModel(item.Id);
						}
					});
				}
			});
		}
		//  查询
	$scope.searchPost = function() {
			getNodeModelList();
		}
		//   回车查询
	$scope.sbgzKeyup = function(e) {
		var keycode = window.event ? e.keyCode : e.which;
		if(keycode == 13) {
			getNodeModelList();
		}
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
				getNodeModelList(pageindex, $scope.page.oneSize);
			}
		}
		/*  -------------------- 分页、页码  -----------------------  */

	$scope.run = function() {
		getNodeModelList();
	}
	$scope.run();
}]);

//设备管理-节点配置-节点型号-添加节点型号弹窗
app.controller("modalJdxhContr", ['$scope', 'httpService', '$modal', '$modalInstance', 'items', 'formValidate', 'toaster', function($scope, httpService, $modal, $modalInstance, items, formValidate, toaster) {
	console.log("设备管理-节点配置-节点型号-添加节点型号弹窗");
	console.log("添加节点型号弹窗", items);
	$scope.items = items;
	$scope.form = {
			"Id": "",
			"Name": "",
			"Description": ""
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

	//节点型号--添加
	var saveNodeModelAdd = function() {
			if(!(formValidate($scope.form.Name).minLength(0).outMsg(2400).isOk)) return false;
			var url = config.HttpUrl + "/device/saveNodeModel";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					Id: getGUIDs(),
					Name: $scope.form.Name,
					Description: $scope.form.Description
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("节点型号添加", data);
				if(data.Rcode == "1000") {
					toaster.pop('success', '添加成功！');
					$modalInstance.close(true);
				} else {
					toaster.pop('danger', data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//节点型号--修改
	var saveNodeModel = function(Id) {
			if(!(formValidate($scope.form.Name).minLength(0).outMsg(2400).isOk)) return false;
			if(!Id) return false;
			var url = config.HttpUrl + "/device/saveNodeModel";
			var data = {
				Auth: {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
				Para: {
					Id: $scope.form.Id,
					Name: $scope.form.Name,
					Description: $scope.form.Description
				}
			}
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				console.log("节点型号修改", data);
				if(data.Rcode == "1000") {
					toaster.pop('success', '修改成功！');
					$modalInstance.close(true);
				} else {
					toaster.pop('danger', data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//  保存按钮
	$scope.ok = function() {

			if($scope.items.operate == "edit") {
				saveNodeModel($scope.items.item.Id);
			} else {
				saveNodeModelAdd();
			}
		}
		//	取消按钮
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	}

	$scope.run = function() {
		if(items.operate == "edit" || items.operate == "see") {
			$scope.form = $.extend({}, $scope.form, $scope.items.item);
		}
	}
	$scope.run();
}]);
