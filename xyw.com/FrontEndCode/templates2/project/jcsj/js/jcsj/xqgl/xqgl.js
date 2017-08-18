'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 基础数据-校区管理
 */

/*    校区管理     */
app.controller("jcsjXqglContr", ['$scope', '$state', 'httpService', '$modal', 'toaster', function($scope, $state, httpService, $modal, toaster) {
	//$state.go("app.qxgl",false);
	console.log("校区管理")

	$scope.xqglData = {
		"campusList": [],
		"buildingList": [],
		"floorsList": []
	}

	//    校区管理 树
	$scope.xqgl_data = [];
	//    校区管理 右边列表
	$scope.xqglRightData = {
		"item": {},
		"list": []
	};

	//   page
	$scope.backPage = {
		PageIndex: 1,
		PageSize: 10
	}

	//   树 已展开 保存数组
	var treeIsOpen = [];

	//    getall
	var getAll = function() {
		//   是展开时加入展开数组
		treeIsOpen = [];
		dg_tree($scope.xqgl_data, function(item) {
			if(item.expanded == true) {
				switch(item.level) {
					case 0:
						//
						break;
					case 1:
						treeIsOpen.push({
							"level": item.level,
							"Campusid": item.Campusid,
							"expanded": true
						});
						break;
					case 2:
						treeIsOpen.push({
							"level": item.level,
							"Buildingid": item.Buildingid,
							"expanded": true
						});
						break;
					case 3:
						treeIsOpen.push({
							"level": item.level,
							"Floorsid": item.Floorsid,
							"expanded": true
						});
						break;
				}
			}
		});
		//
		var url = config.HttpUrl + "/basicset/getall";
		var data = {};
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			console.log("getall", data)
			if(data.Rcode == "1000") {
				$scope.xqglData.campusList = data.Result[0];
				$scope.xqglData.buildingList = data.Result[1];
				$scope.xqglData.floorsList = data.Result[2];

				$scope.xqgl_data = outTree($scope.xqglData.campusList, $scope.xqglData.buildingList, $scope.xqglData.floorsList);
				$scope.xqgl_data[0].expanded = true;
				var objitem = $scope.xqglRightData.item;
				//	默认展开
				if(objitem == null || angular.equals({}, objitem)) {
					$scope.xqgl_data[0].selected = true;
					$scope.xqglRightData.item = $scope.xqgl_data[0];
					campuslist();
				}
				//  加入选中
				dg_tree($scope.xqgl_data, function(item) {
					if($scope.xqglRightData.item.level == item.level && $scope.xqglRightData.item.Campusid == item.Floorsid) {
						item.selected = true;
						//return;
					}
				});
				//   加入展开
				for(var a = 0; a < treeIsOpen.length; a++) {
					if(treeIsOpen[a].level == 1) {
						dg_tree($scope.xqgl_data, function(item) {
							if(item.level == 1 && item.Campusid == treeIsOpen[a].Campusid) {
								item.expanded = true;
								//return;
							}
						});
					}
					if(treeIsOpen[a].level == 2) {
						dg_tree($scope.xqgl_data, function(item) {
							if(item.level == 2 && item.Buildingid == treeIsOpen[a].Buildingid) {
								item.expanded = true;
								//return;
							}
						});
					}
					if(treeIsOpen[a].level == 3) {
						dg_tree($scope.xqgl_data, function(item) {
							if(item.level == 3 && item.Floorsid == treeIsOpen[a].Floorsid) {
								item.expanded = true;
								//return;
							}
						});
					}
				}
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	//  递归树
	var dg_tree = function(tree, fn) {
		for(var i = 0; i < tree.length; i++) {
			fn(tree[i]);
			//console.log(tree[i],tree[i].children)
			if(tree[i].children.length > 0) {
				dg_tree(tree[i].children, fn);
			}
		}
	}

	/*  生成校区楼栋教室树  */
	//   obj,obj,obj
	var outTree = function(campusList, buildingList, floorsList) {
			var tree = [];
			//  计数
			var n1 = 0,
				n2 = 0;
			//   校区
			for(var c in campusList) {
				//  名称
				campusList[c].label = campusList[c].Campusname;
				campusList[c].level = 1;
				tree.push(campusList[c]);
				n2 = 0;
				//   楼栋
				for(var b in buildingList) {
					if(!("children" in tree[n1])) tree[n1].children = [];
					if(tree[n1].Campusid == buildingList[b].Campusid) {
						//  名称
						buildingList[b].label = buildingList[b].Buildingname;
						buildingList[b].level = 2;
						//   插入数组
						tree[n1].children.push(buildingList[b]);
						//   楼层
						for(var f in floorsList) {
							if(!("children" in tree[n1].children[n2])) tree[n1].children[n2].children = [];
							//    加上children属性
							if(!("children" in floorsList[f])) floorsList[f].children = [];

							if(tree[n1].children[n2].Buildingid == floorsList[f].Buildingid) {
								//  名称
								floorsList[f].label = floorsList[f].Floorname;
								floorsList[f].level = 3;
								//   插入数组
								tree[n1].children[n2].children.push(floorsList[f]);
							}
						}
						n2++;
					}
				}
				n1++;
			}
			return [{
				label: '全部校区',
				children: tree,
				level: 0
			}];
		}
		//////////////////////////////////////////////////////////////////////
		//    校区查询
	var campuslist = function() {
		var url = config.HttpUrl + "/system/bs/campuslist";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("校区查询", data)
			if(data.Rcode == "1000") {
				$scope.xqglRightData.list = data.Result.PageData;
				//   分页
				var objPage = {
					PageCount: data.Result.PageCount,
					PageIndex: data.Result.PageIndex,
					PageSize: data.Result.PageSize,
					RecordCount: data.Result.PageCount
				};
				if((objPage.RecordCount % objPage.PageSize) == 0) {
					objPage.PageCount = (objPage.RecordCount / objPage.PageSize);
				} else {
					objPage.PageCount = parseInt((objPage.RecordCount / objPage.PageSize)) + 1;
				}
				$scope.backPage = pageFn(objPage, 5);
			} else {
				console.log(data.Reason);
				$scope.xqglRightData.list = [];
			}
		}, function(reason) {}, function(update) {});
	};

	//    删除校区
	var campusdel = function(Id) {
		if(!Id) return false;
		var url = config.HttpUrl + "/system/bs/campusdel";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Id": Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("删除校区", data)
			if(data.Rcode == "1000") {
				campuslist();
				getAll();
				toaster.pop('success', '删除成功！');
			} else {
				toaster.pop('success', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	////////////////////////////////////////////////////////////////////////////
	//    楼栋查询
	var buildinglist = function() {
		var url = config.HttpUrl + "/system/bs/buildinglist";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Campusid": $scope.xqglRightData.item.Campusid,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("楼栋查询", data)
			if(data.Rcode == "1000") {
				if(data.Result != null) {
					$scope.xqglRightData.list = data.Result.PageData;
					//   分页
					var objPage = {
						PageCount: data.Result.PageCount,
						PageIndex: data.Result.PageIndex,
						PageSize: data.Result.PageSize,
						RecordCount: data.Result.PageCount
					};
					if((objPage.RecordCount % objPage.PageSize) == 0) {
						objPage.PageCount = (objPage.RecordCount / objPage.PageSize);
					} else {
						objPage.PageCount = parseInt((objPage.RecordCount / objPage.PageSize)) + 1;
					}
					$scope.backPage = pageFn(objPage, 5);
				} else {
					$scope.xqglRightData.list = [];
				}
			} else {
				console.log(data.Reason);
				$scope.xqglRightData.list = [];
			}
		}, function(reason) {}, function(update) {});
	};

	//    删除楼栋
	var buildingdel = function(Id) {
		if(!Id) return false;
		var url = config.HttpUrl + "/system/bs/buildingdel";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Id": Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("删除楼栋", data)
			if(data.Rcode == "1000") {
				buildinglist();
				getAll();
				toaster.pop('success', '删除成功！');
			} else {
				toaster.pop('success', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	//////////////////////////////////////////////////////////////////////////
	//    楼层查询
	var floorslist = function() {
		var url = config.HttpUrl + "/system/bs/floorslist";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Buildingid": $scope.xqglRightData.item.Buildingid,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("楼层查询", data)
			if(data.Rcode == "1000") {
				$scope.xqglRightData.list = data.Result.PageData;
				//   分页
				var objPage = {
					PageCount: data.Result.PageCount,
					PageIndex: data.Result.PageIndex,
					PageSize: data.Result.PageSize,
					RecordCount: data.Result.PageCount
				};
				if((objPage.RecordCount % objPage.PageSize) == 0) {
					objPage.PageCount = (objPage.RecordCount / objPage.PageSize);
				} else {
					objPage.PageCount = parseInt((objPage.RecordCount / objPage.PageSize)) + 1;
				}
				$scope.backPage = pageFn(objPage, 5);
			} else {
				console.log(data.Reason);
				$scope.xqglRightData.list = [];
			}
		}, function(reason) {}, function(update) {});
	};

	//    删除楼层
	var floorsdel = function(Id) {
		var url = config.HttpUrl + "/system/bs/floorsdel";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Id": Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("删除楼层", data)
			if(data.Rcode == "1000") {
				floorslist();
				getAll();
				toaster.pop('success', '删除成功！');
			} else {
				toaster.pop('success', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	////////////////////////////////////////////////////////////////////////////

	//    教室查询
	var classroomslist = function() {
		var url = config.HttpUrl + "/system/bs/classroomslist";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Floorsid": $scope.xqglRightData.item.Floorsid,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("教室查询", data)
			if(data.Rcode == "1000") {
				$scope.xqglRightData.list = data.Result.PageData;
				//   分页
				var objPage = {
					PageCount: data.Result.PageCount,
					PageIndex: data.Result.PageIndex,
					PageSize: data.Result.PageSize,
					RecordCount: data.Result.PageCount
				};
				if((objPage.RecordCount % objPage.PageSize) == 0) {
					objPage.PageCount = (objPage.RecordCount / objPage.PageSize);
				} else {
					objPage.PageCount = parseInt((objPage.RecordCount / objPage.PageSize)) + 1;
				}
				$scope.backPage = pageFn(objPage, 5);
			} else {
				console.log(data.Reason);
				$scope.xqglRightData.list = [];
			}
		}, function(reason) {}, function(update) {});
	};

	//    教室删除
	var classroomsdel = function(Id) {
		var url = config.HttpUrl + "/system/bs/classroomsdel";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Id": Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("教室删除", data)
			if(data.Rcode == "1000") {
				classroomslist();
				getAll();
				toaster.pop('success', '删除成功！');
			} else {
				toaster.pop('success', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	/////////////////////////////////////////////////////////////////////////////

	//添加按钮功能
	$scope.openModalAdd = function(str, active, item) {
		if(!("level" in active)) {
			$modal.open({
				templateUrl: 'modal/modal_alert_all.html',
				controller: 'modalAlert2Conter',
				resolve: {
					items: function() {
						return {
							"type": 'warning',
							"msg": '请选择选项'
						};
					}
				}
			});
		}
		var modalInstance = $modal.open({
			templateUrl: '../project/jcsj/html/jcsj/xqgl/modal_add.html',
			controller: 'modalXqglContr',
			windowClass: 'm-modal-xqgl',
			resolve: {
				items: function() {
					return {
						'str': str,
						"active": active,
						"item": item
					};
				}
			}
		});

		modalInstance.result.then(function(bul) {
			console.log(bul)
			if(bul) {
				//   刷新列表
				switch($scope.xqglRightData.item.level) {
					case 0:
						campuslist();
						break;
					case 1:
						buildinglist();
						break;
					case 2:
						floorslist();
						break;
					case 3:
						classroomslist();
						break;
				}
				getAll();
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    树点击
	$scope.xqgl_tree_handler = function(branch) {

		//
		console.log(branch)
		$scope.backPage.PageIndex = 1;
		$scope.backPage.PageSize = 10;
		$scope.xqglRightData.item = branch;
		switch($scope.xqglRightData.item.level) {
			case 0:
				campuslist();
				break;
			case 1:
				buildinglist();
				break;
			case 2:
				floorslist();
				break;
			case 3:
				classroomslist();
				break;
		}
		if($scope.xqglRightData.item.level != 0) {
			$scope.xqgl_data[0].selected = false;
		}
	};
	
	
	//  删除前提示
	var deleteBefore = function(Id,Code,delFn) {
//		var url = config.HttpUrl + "/device/onDeletingDeviceModelStatusCmd";
//		var data = {
//			Auth: {
//				"Usersid": config.GetUser().Usersid,
//				"Rolestype": config.GetUser().Rolestype,
//				"Token": config.GetUser().Token,
//				"Os": "WEB"
//			},
//			Para: {
//				ModelId: Id,
//				StatusCode: Code
//			}
//		}
//		var promise = httpService.ajaxPost(url, data);
//		promise.then(function(data) {
//			console.log("删除前提示", data);
//			if(data.Rcode == "1000") {
//				delFn(data.Result);
//			} else {
//				return false;
//			}
//		});

		delFn(false);

	}
	
	

	//   删除弹窗
	$scope.deleteItem = function(active, item) {
		//
		var msg_xq = [];
		switch(active.level) {
			case 0:
				//   删除学院
				msg_xq = ['确定删除当前校区吗？','警告：此操作将删除当前校区下所有教学楼、楼层、教室，同时会影响关联当前校区的设备、节点、场景、日志，故障、统计分析、预警数据完整性！'];
				break;
			case 1:
				msg_xq = ['确定删除当前教学楼吗？','警告：此操作将删除当前教学楼下所有楼层、教室，同时会影响关联当前教学楼的设备、节点、场景、日志，故障、统计分析、预警数据完整性！'];
				break;
			case 2:
				msg_xq = ['确定删除当前楼层吗？','警告：此操作将删除当前楼层下所有教室，同时会影响关联当前楼层的设备、节点、场景、日志，故障、统计分析、预警数据完整性！'];
				break;
			case 3:
				msg_xq = ['确定删除当前教室吗？','警告：此操作将删除当前教室下所有教室，同时会影响关联当前教室的设备、节点、场景、日志，故障、统计分析、预警数据完整性！'];
				break;
		}
		
		//
		var modalInstance = $modal.open({
			templateUrl: 'modal/modal_alert_all.html',
			controller: 'modalAlert2Conter',
			resolve: {
				items: function() {
					return {
						"type": 'warning',
						"msg": msg_xq[0]
					};
				}
			}
		});
		modalInstance.result.then(function(bul) {
			console.log("bul", bul);
			if(bul) {
				deleteBefore(item.Id,'code', function(rBul) {
					//   是否有关联内容
					if(rBul){
						var modalInstance2 = $modal.open({
							templateUrl: 'modal/modal_alert_all.html',
							controller: 'modalAlert2Conter',
							resolve: {
								items: function() {
									return {
										"type": 'warning',
										"msg": msg_xq[1]
									};
								}
							}
						});
						modalInstance2.result.then(function(bul) {
							console.log("bul", bul);
							if(bul) {
								//    删除
								switch(active.level) {
									case 0:
										//   删除学院
										campusdel(item.Campusid);
										break;
									case 1:
										buildingdel(item.Buildingid);
										break;
									case 2:
										floorsdel(item.Floorsid);
										break;
									case 3:
										classroomsdel(item.Classroomid);
										break;
								}
							}
						});
					}else{
						//  删除
						switch(active.level) {
							case 0:
								//   删除学院
								campusdel(item.Campusid);
								break;
							case 1:
								buildingdel(item.Buildingid);
								break;
							case 2:
								floorsdel(item.Floorsid);
								break;
							case 3:
								classroomsdel(item.Classroomid);
								break;
						}
					}
				});
			}
		}); 
		
		

	}

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
	};
	//  翻页
	$scope.pageClick = function(pageindex) {
		if(!(Number(pageindex) > 0)) return false;
		if(pageindex > 0 && pageindex <= $scope.backPage.PageCount) {
			$scope.backPage.PageIndex = pageindex;
			switch($scope.xqglRightData.item.level) {
				case 0:
					campuslist();
					break;
				case 1:
					buildinglist();
					break;
				case 2:
					floorslist();
					break;
				case 3:
					classroomslist();
					break;
			}
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */

	$scope.run = function() {
		getAll();
	}
	$scope.run();

}]);

app.controller("modalXqglContr", ['$scope', 'httpService', '$modal', '$modalInstance', 'items', 'formValidate', 'toaster', function($scope, httpService, $modal, $modalInstance, items, formValidate, toaster) {
	console.log("基础数据-校区管理-弹窗");

	$scope.items = items;
	console.log(items);

	$scope.form = {
		//   校区
		"campus": {},
		//   楼栋
		"building": {},
		//   楼层
		"floors": {},
		//   教室
		"classrooms": {}
	}

	$scope.form.campus = {
		"Campuscode": "",
		"Campusname": "",
		"Campusicon": "",
		"Id": null
	}

	$scope.form.building = {
		"Buildingcode": "",
		"Buildingname": "",
		"Buildingicon": "",
		"Campusid": null,
		"Id": null
	}
	$scope.form.floors = {
		"Floorscode": "",
		"Floorname": "",
		"FloorsImage": "",
		"Buildingid": null,
		"Id": null
	}
	$scope.form.classrooms = {
		"Classroomscode": "",
		"Classroomsname": "",
		"Classroomicon": "",
		"Floorsid": null,
		"Seatsnumbers": null,
		"Classroomstype": "普通",
		//  ----
		"ClassroomstypeItem": {
			'val': '普通',
			'title': '普通'
		},
		//  ----
		"ClassroomstypeItems": [{
			'val': '普通',
			'title': '普通'
		}, {
			'val': '多功能',
			'title': '多功能'
		}],
		"Classroomstate": 0,
		"Notes": "",
		"Maxy": null,
		"Miny": null,
		"Maxx": null,
		"Minx": null,
		"Id": null
	}

	//   上传图片
	$scope.upimglist = [];

	//    清除图片
	$scope.closePic = function(index) {
		$scope.upimglist.splice(index, 1);
	}

	//   select  教室属性
	$scope.changeClassroomstypeItem = function(item) {
		$scope.form.classrooms.ClassroomstypeItem = item;
		$scope.form.classrooms.Classroomstype = item.val;
	}

	//    验证校区代码
	$scope.changeSubjectcode = function() {
		//   本级学院增加长度
		var xk_lang = 4;
		//   学院代码
		var code = '';
		//   
		if($scope.items.active.level == 0) {
			code = $scope.form.campus.Campuscode;
		} else if($scope.items.active.level == 1) {
			code = $scope.form.building.Buildingcode;
		} else if($scope.items.active.level == 2) {
			code = $scope.form.floors.Floorscode;
		} else if($scope.items.active.level == 3) {
			code = $scope.form.classrooms.Classroomscode;
		}

		//   上级学院代码
		var pcode = "";
		if($scope.items.active.level == 0) {
			pcode = "";
		} else if($scope.items.active.level == 1) {
			pcode = $scope.items.active.Campuscode;
		} else if($scope.items.active.level == 2) {
			pcode = $scope.items.active.Buildingcode;
		} else if($scope.items.active.level == 3) {
			pcode = $scope.items.active.Floorscode;
		}

		//    验证输入必须为数字
		if(!(/^[a-zA-Z0-9-]*$/.test(code))) {

			$modal.open({
				templateUrl: 'modal/modal_alert_all.html',
				controller: 'modalAlert2Conter',
				resolve: {
					items: function() {
						return {
							"type": 'danger',
							"msg": '只能输入上级代码加4位的数字或字母！'
						};
					}
				}
			});
			//alert('只能输入上级代码加4位的数字！');
			code = pcode;
		}

		//
		if(code.length < pcode.length) {
			//   里面没有上级学院代码  则 放入上级学院代码
			code = pcode;
		} else {
			//   有  验证前面的数字是不是上级学院代码
			if(code.substr(0, pcode.length) == pcode) {
				//   是上级学院代码 。限制输入长度
				if(code.length > pcode.length + xk_lang) {
					code = code.substr(0, code.length - 1);
				} else {
					console.log(code);
				}
			} else {
				//   不是上级学院代码 。 放入上级学院代码
				code = pcode;
			}
		}
		//
		//$scope.form.subject.Subjectcode = code;
		if($scope.items.active.level == 0) {
			$scope.form.campus.Campuscode = code;
		} else if($scope.items.active.level == 1) {
			$scope.form.building.Buildingcode = code;
		} else if($scope.items.active.level == 2) {
			$scope.form.floors.Floorscode = code;
		} else if($scope.items.active.level == 3) {
			$scope.form.classrooms.Classroomscode = code;
		}
	}

	//    验证提交
	var xq_code = function() {
		//
		var temp_bol = false;
		switch($scope.items.active.level) {
			case 0:
				if($scope.form.campus.Campuscode == "" || $scope.form.campus.Campuscode == undefined) {
					temp_bol = true;
				}
				break;
			case 1:
				if($scope.form.building.Buildingcode == $scope.items.active.Campuscode) {
					temp_bol = true;
				}
				break;
			case 2:
				if($scope.form.floors.Floorscode == $scope.items.active.Buildingcode) {
					temp_bol = true;
				}
				break;
			case 2:
				if($scope.form.classrooms.Classroomscode == $scope.items.active.Floorscode) {
					temp_bol = true;
				}
				break;
		}
		if(temp_bol) {
			$modal.open({
				templateUrl: 'modal/modal_alert_all.html',
				controller: 'modalAlert2Conter',
				resolve: {
					items: function() {
						return {
							"type": 'danger',
							"msg": '代码不能与上级代码相同！'
						};
					}
				}
			});
			return false;
		} else {
			return true;
		}
	}

	//   校区添加
	var campusadd = function() {
		if(!(formValidate($scope.form.campus.Campuscode).minLength(0).outMsg(2803).isOk)) return false;
		if(!(formValidate($scope.form.campus.Campusname).minLength(0).outMsg(2804).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/campusadd";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Campuscode": $scope.form.campus.Campuscode,
			"Campusname": $scope.form.campus.Campusname,
			"Campusicon": $scope.form.campus.Campusicon
		};

		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("校区添加", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '添加成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	//   校区修改
	var campuschange = function() {
		if(!(formValidate($scope.form.campus.Campuscode).minLength(0).outMsg(2803).isOk)) return false;
		if(!(formValidate($scope.form.campus.Campusname).minLength(0).outMsg(2804).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/campuschange";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Campuscode": $scope.form.campus.Campuscode,
			"Campusname": $scope.form.campus.Campusname,
			"Campusicon": $scope.form.campus.Campusicon,
			"Id": $scope.form.campus.Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("校区修改", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '修改成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	/////////////////////////////
	//   楼栋添加
	var buildingadd = function() {
		if(!(formValidate($scope.form.building.Buildingcode).minLength(0).outMsg(2805).isOk)) return false;
		if(!(formValidate($scope.form.building.Buildingname).minLength(0).outMsg(2806).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/buildingadd";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Buildingcode": $scope.form.building.Buildingcode,
			"Buildingname": $scope.form.building.Buildingname,
			"Buildingicon": $scope.form.building.Buildingicon,
			"Campusid": $scope.form.building.Campusid
		};

		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("楼栋添加", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '添加成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	//   楼栋修改
	var buildingchange = function() {
		if(!(formValidate($scope.form.building.Buildingcode).minLength(0).outMsg(2805).isOk)) return false;
		if(!(formValidate($scope.form.building.Buildingname).minLength(0).outMsg(2806).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/buildingchange";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Buildingcode": $scope.form.building.Buildingcode,
			"Buildingname": $scope.form.building.Buildingname,
			"Buildingicon": $scope.form.building.Buildingicon,
			"Campusid": $scope.form.building.Campusid,
			"Id": $scope.form.building.Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("楼栋修改", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '修改成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	/////////////////////////////////////////////

	//   楼层添加
	var floorsadd = function() {
		if(!(formValidate($scope.form.floors.Floorscode).minLength(0).outMsg(2807).isOk)) return false;
		if(!(formValidate($scope.form.floors.Floorname).minLength(0).outMsg(2808).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/floorsadd";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Floorscode": $scope.form.floors.Floorscode,
			"Floorname": $scope.form.floors.Floorname,
			"FloorsImage": $scope.form.floors.FloorsImage,
			"Buildingid": $scope.form.floors.Buildingid
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("楼层添加", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '添加成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	//   楼层修改
	var floorschange = function() {
		if(!(formValidate($scope.form.floors.Floorscode).minLength(0).outMsg(2807).isOk)) return false;
		if(!(formValidate($scope.form.floors.Floorname).minLength(0).outMsg(2808).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/floorschange";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Floorscode": $scope.form.floors.Floorscode,
			"Floorname": $scope.form.floors.Floorname,
			"FloorsImage": $scope.form.floors.FloorsImage,
			"Buildingid": $scope.form.floors.Buildingid,
			"Id": $scope.form.floors.Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("楼层修改", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '修改成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	/////////////////////////////////////////////////////////////////

	//   教室添加
	var classroomsadd = function() {
		if(!(formValidate($scope.form.classrooms.Classroomscode).minLength(0).outMsg(2809).isOk)) return false;
		if(!(formValidate($scope.form.classrooms.Classroomsname).minLength(0).outMsg(2810).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/classroomsadd";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Classroomsname": $scope.form.classrooms.Classroomsname,
			"Classroomscode": $scope.form.classrooms.Classroomscode,
			"Classroomicon": $scope.form.classrooms.Classroomicon,
			"Floorsid": $scope.form.classrooms.Floorsid,
			"Seatsnumbers": Number($scope.form.classrooms.Seatsnumbers),
			"Classroomstype": $scope.form.classrooms.Classroomstype,
			"Classroomstate": $scope.form.classrooms.Classroomstate,
			"Notes": $scope.form.classrooms.Notes
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("教室添加", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '添加成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	//   教室修改
	var classroomschange = function() {
		if(!(formValidate($scope.form.classrooms.Classroomscode).minLength(0).outMsg(2809).isOk)) return false;
		if(!(formValidate($scope.form.classrooms.Classroomsname).minLength(0).outMsg(2810).isOk)) return false;

		//   验证提交
		if(!xq_code()) return false;
		//
		$scope.changeSubjectcode();

		var url = config.HttpUrl + "/system/bs/classroomschange";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Classroomsname": $scope.form.classrooms.Classroomsname,
			"Classroomscode": $scope.form.classrooms.Classroomscode,
			"Classroomicon": $scope.form.classrooms.Classroomicon,
			"Floorsid": $scope.form.classrooms.Floorsid,
			"Seatsnumbers": Number($scope.form.classrooms.Seatsnumbers),
			"Classroomstype": $scope.form.classrooms.Classroomstype,
			"Classroomstate": $scope.form.classrooms.Classroomstate,
			"Notes": $scope.form.classrooms.Notes,
			"Id": $scope.form.classrooms.Id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("教室修改", data)
			if(data.Rcode == "1000") {
				toaster.pop('success', '修改成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	///////////////////////////////////////////////////////////////////////

	//
	$scope.ok = function() {
		if($scope.items.str == "add") {
			switch($scope.items.active.level) {
				case 0:
					if($scope.upimglist.length > 0) {
						$scope.form.campus.Campusicon = $scope.upimglist[0].Result;
					} else {
						$scope.form.campus.Campusicon = "";
					}
					campusadd();
					break;
				case 1:
					if($scope.upimglist.length > 0) {
						$scope.form.building.Buildingicon = $scope.upimglist[0].Result;
					} else {
						$scope.form.building.Buildingicon = "";
					}
					$scope.form.building.Campusid = $scope.items.active.Campusid;
					buildingadd();
					break;
				case 2:
					if($scope.upimglist.length > 0) {
						$scope.form.floors.FloorsImage = $scope.upimglist[0].Result;
					} else {
						$scope.form.floors.FloorsImage = "";
					}
					$scope.form.floors.Buildingid = $scope.items.active.Buildingid;
					floorsadd();
					break;
				case 3:
					if($scope.upimglist.length > 0) {
						$scope.form.classrooms.Classroomicon = $scope.upimglist[0].Result;
					} else {
						$scope.form.classrooms.Classroomicon = "";
					}
					$scope.form.classrooms.Floorsid = $scope.items.active.Floorsid;
					classroomsadd();
					break;
			}
		}
		if($scope.items.str == "edit") {
			switch($scope.items.active.level) {
				case 0:
					if($scope.upimglist.length > 0) {
						$scope.form.campus.Campusicon = $scope.upimglist[0].Result;
					} else {
						$scope.form.campus.Campusicon = "";
					}
					campuschange();
					break;
				case 1:
					if($scope.upimglist.length > 0) {
						$scope.form.building.Buildingicon = $scope.upimglist[0].Result;
					} else {
						$scope.form.building.Buildingicon = "";
					}
					//$scope.form.building.Campusid = $scope.items.active.Campusid;
					buildingchange();
					break;
				case 2:
					if($scope.upimglist.length > 0) {
						$scope.form.floors.FloorsImage = $scope.upimglist[0].Result;
					} else {
						$scope.form.floors.FloorsImage = "";
					}
					//$scope.form.floors.Buildingid = $scope.items.active.Buildingid;
					floorschange();
					break;
				case 3:
					if($scope.upimglist.length > 0) {
						$scope.form.classrooms.Classroomicon = $scope.upimglist[0].Result;
					} else {
						$scope.form.classrooms.Classroomicon = "";
					}
					//$scope.form.classrooms.Floorsid = $scope.items.active.Floorsid;
					classroomschange();
					break;
			}
		}
	}

	//	取消按钮
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	}

	$scope.run = function() {
		if($scope.items.str == "add") {
			//
			switch($scope.items.active.level) {
				case 0:
					$scope.form.campus.Campuscode = "";
					break;
				case 1:
					$scope.form.building.Buildingcode = $scope.items.active.Campuscode;
					break;
				case 2:
					$scope.form.floors.Floorscode = $scope.items.active.Buildingcode;
					break;
				case 3:
					$scope.form.classrooms.Classroomscode = $scope.items.active.Floorscode;
					break;
			}
		}
		if($scope.items.str == "edit") {
			switch($scope.items.active.level) {
				case 0:
					$scope.form.campus.Campuscode = $scope.items.item.Campuscode;
					$scope.form.campus.Campusname = $scope.items.item.Campusname;
					$scope.form.campus.Campusicon = $scope.items.item.Campusicon;
					$scope.form.campus.Id = $scope.items.item.Campusid;
					if($scope.form.campus.Campusicon.length > 0) $scope.upimglist[0] = {
						Result: $scope.form.campus.Campusicon
					};
					break;
				case 1:
					$scope.form.building = $.extend({}, $scope.form.building, $scope.items.item);
					$scope.form.building.Id = $scope.items.item.Buildingid;
					if($scope.form.building.Buildingicon.length > 0) $scope.upimglist[0] = {
						Result: $scope.form.building.Buildingicon
					};
					break;
				case 2:
					$scope.form.floors = $.extend({}, $scope.form.floors, $scope.items.item);
					$scope.form.floors.Id = $scope.items.item.Floorsid;
					if($scope.form.floors.FloorsImage.length > 0) $scope.upimglist[0] = {
						Result: $scope.form.floors.FloorsImage
					};
					break;
				case 3:
					$scope.form.classrooms = $.extend({}, $scope.form.classrooms, $scope.items.item);
					$scope.form.classrooms.Id = $scope.items.item.Classroomid;
					if($scope.form.classrooms.Classroomicon.length > 0) $scope.upimglist[0] = {
						Result: $scope.form.classrooms.Classroomicon
					};
					break;
			}
		}
	}
	$scope.run();

}]);