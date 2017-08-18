'use strict';

/**
 * 设备管理路由
 * 
 */

app.config(
		['$stateProvider', '$urlRouterProvider',
			function($stateProvider, $urlRouterProvider) {
				$stateProvider
					/*  //////////////////  设备管理   /////////////////////   */
					//   设备管理
						.state('app.sbgl', {
							url: '/sbgl',
							templateUrl: '../project/sbgl/html/sbgl/index.html',
							controller: 'sbglindexContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../vendor/echarts/echarts.min.js',
													'../project/sbgl/js/sbgl/sbgl.js'
												]
										);
									}
								]
							}
						})
						//  设备型号管理
						.state('app.xhgl', {
							url: '/sbgl/xhgl',
							templateUrl: '../project/sbgl/html/sbgl/sbgl_xhgl.html',
							controller: 'sbglXhglindexContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												['../js/controllers/sbgl.js']
										);
									}
								]
							}
						})
						//   设备日志
						.state('app.sbgl.sbrz', {
							url: '/sbrz',
							templateUrl: '../project/sbgl/html/sbgl/sbrz/index.html',
							controller: 'sbglSbrzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbrz/sbrz.css',
													'../project/sbgl/js/sbgl/sbrz/sbrz.js'
												]
										);
									}
								]
							}
						})


						//   设备故障
						.state('app.sbgl.sbgz', {
							url: '/sbgz',
							templateUrl: '../project/sbgl/html/sbgl/sbgz/index.html',
							controller: 'sbglSbgzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbgz/sbgz.css',
													'../project/sbgl/js/sbgl/sbgz/sbgz.js'
												]
										);
									}
								]
							}
						})

						/*     设备监控         */
						.state('app.sbgl.sbjk', {
							url: '/sbjk',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/index.html',
							controller: 'sbjkContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbjk/sbjk.css',
													'../vendor/jquery/blockUI/jquery.blockUI.js',
													'../project/sbgl/js/sbgl/sbjk/sbjk.js'
												]
										);
									}
								]
							}
						})
						/*     设备监控列表页 教室设备列表         */
						.state('app.sbgl.sbjk.list', {
							url: '/list',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/list.html',
							controller: 'sbjkListContr',
							params:{classroomid:null},
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/js/sbgl/sbjk/sbjk.js'
												]
										);
									}
								]
							}
						})

						/*     设备监控列表页 教室设备控制详情         */
						.state('app.sbgl.sbjk.list.article', {
							url: '/article',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/article.html',
							controller: 'sbjkArticleContr',
							params:{DeviceId:null,page:null}
						})

						/*     设备监控列表 查课表         */
						.state('app.sbgl.sbjk.list.ckb', {
							url: '/list',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/schedule.html',
							controller: 'sbjkListCkbContr',
							params:{classroomid:null},
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/js/sbgl/sbjk/sbjk.js'
												]
										);
									}
								]
							}
						})
						/*  /////////////////// 设备监控 End  //////////////////  */
						//   设备预警
						.state('app.sbgl.sbyj', {
							url: '/sbyj',
							templateUrl: '../project/sbgl/html/sbgl/sbyj/index.html',
							controller: 'sbglSbyjContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbyj/sbyj.css',
													'../project/sbgl/js/sbgl/sbyj/sbyj.js'
												]
										);
									}
								]
							}
						})
						//   联动场景
						.state('app.sbgl.ldcj', {
							url: '/ldcj',
							templateUrl: '../project/sbgl/html/sbgl/ldcj/index.html',
							controller: 'sbglLdcjContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/ldcj/ldcj.css',
													'../project/sbgl/js/sbgl/ldcj/ldcj.js'
												]
										);
									}
								]
							}
						})
						//   设备分析
						.state('app.sbgl.sbfx', {
							url: '/sbfx',
							templateUrl: '../project/sbgl/html/sbgl/sbfx/index.html',
							controller: 'sbglSbfxContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbfx/sbfx.css',
													'../project/sbgl/js/sbgl/sbfx/sbfx.js'
												]
										);
									}
								]
							}
						})

						//   节点配置
						.state('app.sbgl.jdpz', {
							url: '/jdpz',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/index.html',
							controller: 'sbglJdpzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/jdpz.js'
												]
										);
									}
								]
							}
						})

						//   节点配置-节点型号
						.state('app.sbgl.jdpz.jdxh', {
							url: '/jdxh',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/jdxh/index.html',
							controller: 'sbglJdxhContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/jdxh/jdxh.js'
												]
										);
									}
								]
							}
						})

						//   节点配置-控制命令
						.state('app.sbgl.jdpz.kzml', {
							url: '/kzml',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/kzml/index.html',
							controller: 'sbglJdpzKzmlContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/kzml/kzml.js'
												]
										);
									}
								]
							}
						})

						//   节点配置-节点
						.state('app.sbgl.jdpz.jd', {
							url: '/jd',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/jd/index.html',
							controller: 'sbglJdContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/jd/jd.js'
												]
										);
									}
								]
							}
						})


						// --------------  设备配置----------------
						.state('app.sbgl.sbpz', {
							url: '/sbpz',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/index.html',
							controller: 'sbglSbpzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/sbpz.js'
												]
										);
									}
								]
							}
						})

						// --------------  设备配置--故障分类----------------
						.state('app.sbgl.sbpz.gzfl', {
							url: '/gzfl',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/gzfl/index.html',
							controller: 'sbglGzflContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function( $ocLazyLoad ){
										return $ocLazyLoad.load('angularBootstrapNavTree').then(
												function(){
													return $ocLazyLoad.load([
														'../project/sbgl/css/sbgl/sbpz/sbpz.css',
														'../project/sbgl/js/sbgl/sbpz/gzfl/gzfl.js'
													]);
												}
										);
									}]
							}
						})

						// --------------  设备配置--设备型号----------------
						.state('app.sbgl.sbpz.sbxh', {
							url: '/sbxh',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/sbxh/index.html',
//					controller-----控制器
							controller: 'sbglSbxhContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function( $ocLazyLoad ){
										return $ocLazyLoad.load('angularBootstrapNavTree').then(
												function(){
													return $ocLazyLoad.load([
														'../project/sbgl/css/sbgl/sbpz/sbpz.css',
														'../project/sbgl/js/sbgl/sbpz/sbxh/sbxh.js'
													]);
												}
										);
									}]
							}
						})
						// --------------  设备配置--控制命令----------------
						.state('app.sbgl.sbpz.kzml', {
							url: '/kzml',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/kzml/index.html',
//					controller-----控制器
							controller: 'sbglSbpzKzmlContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/kzml/kzml.js'
												]
										);
									}
								]
							}
						})

						// --------------  设备配置--状态命令----------------
						.state('app.sbgl.sbpz.ztml', {
							url: '/ztml',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/ztml/index.html',
//					controller-----控制器
							controller: 'sbglZtmlContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/ztml/ztml.js'
												]
										);
									}
								]
							}
						})
						// --------------  设备配置--故障现象----------------
						.state('app.sbgl.sbpz.gzxx', {
							url: '/gzxx',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/gzxx/index.html',
//					controller-----控制器
							controller: 'sbglGzxxContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function( $ocLazyLoad ){
										return $ocLazyLoad.load('angularBootstrapNavTree').then(
												function(){
													return $ocLazyLoad.load([
														'../project/sbgl/css/sbgl/sbpz/sbpz.css',
														'../project/sbgl/js/sbgl/sbpz/gzxx/gzxx.js'
													]);
												}
										);
									}]
							}
						})
						/*   --------设备配置--状态命令-- 添加 code     */
						.state('app.sbgl.sbpz.ztml.code', {
							url: '/code',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/ztml/code/index.html',
							controller: 'modalZtmlAddContr',
							params: {
								mid:"",
								code: "",
								code_name:""
							},
							resolve: {
								deps: ['$ocLazyLoad',
									function ($ocLazyLoad) {
										return $ocLazyLoad.load(['ui.select']).then(
												function () {
													return $ocLazyLoad.load(
															[
																'../project/sbgl/css/sbgl/sbpz/sbpz.css',
																'../project/sbgl/js/sbgl/sbpz/ztml/ztml.js'
															]
													);
												}
										);
									}
								]
							}
						})
						// --------------  设备配置--设备管理者----------------
						.state('app.sbgl.sbpz.sb', {
							url: '/sbglz',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/sbglz/index.html',
//					controller-----控制器
							controller: 'sbglSbglzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/sbglz/sbglz.js'
												]
										);
									}
								]
							}
						})
						/*  //////////////////  设备管理 End /////////////////////   */


				/*  //////////////////  基础数据  ///////////////////////  */
				/*   基础数据       */
				.state('app.jcsj', {
					url: '/jcsj',
					templateUrl: '../project/jcsj/html/jcsj/index.html',
					controller: 'jcsjContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/jcsj.js'
									]
								);
							}
						]
					}
				})
				
				/*   校区管理       */
				.state('app.jcsj.xqgl', {
					url: '/xqgl',
					templateUrl: '../project/jcsj/html/jcsj/xqgl/index.html',
					controller: 'jcsjXqglContr',
					resolve: {
						deps: ['$ocLazyLoad',
                        function( $ocLazyLoad ){
                        	return $ocLazyLoad.load('angularBootstrapNavTree').then(
                            	function(){
                                	return $ocLazyLoad.load([
                                		'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/xqgl/xqgl.js'
                                	]);
                            	} 
                        	);
                    	}]
					}
				})
				
				/*   学院管理       */
				.state('app.jcsj.xygl', {
					url: '/xygl',
					templateUrl: '../project/jcsj/html/jcsj/xygl/index.html',
					controller: 'jcsjXyglContr',
					resolve: {
						deps: ['$ocLazyLoad',
                        function( $ocLazyLoad ){
                        	return $ocLazyLoad.load('angularBootstrapNavTree').then(
                            	function(){
                                	return $ocLazyLoad.load([
                                		'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/xygl/xygl.js'
                                	]);
                            	} 
                        	);
                    	}]
					}
				})
				
				/*   学科管理       */
				.state('app.jcsj.xkgl', {
					url: '/xkgl',
					templateUrl: '../project/jcsj/html/jcsj/xkgl/index.html',
					controller: 'jcsjXkglContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function( $ocLazyLoad ){
                        	return $ocLazyLoad.load('angularBootstrapNavTree').then(
                            	function(){
                                	return $ocLazyLoad.load([
                                		'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/xkgl/xkgl.js'
                                	]);
                            	} 
                        	);
                    	}]
					}
				})
				
				
				/*   用户管理       */
				.state('app.jcsj.yhgl', {
					url: '/yhgl',
					templateUrl: '../project/jcsj/html/jcsj/yhgl/index.html',
					controller: 'jcsjYhglContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/yhgl/yhgl.js'
									]
								);
							}
						]
					}
				})
				
				/*   职工管理       */
				.state('app.jcsj.yhgl.zggl', {
					url: '/zggl',
					templateUrl: '../project/jcsj/html/jcsj/yhgl/zggl/index.html',
					controller: 'jcsjYhglZgglContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/yhgl/zggl/zggl.js'
									]
								);
							}
						]
					}
				})
				/*   教师管理       */
				.state('app.jcsj.yhgl.jsgl', {
					url: '/jsgl',
					templateUrl: '../project/jcsj/html/jcsj/yhgl/jsgl/index.html',
					controller: 'jcsjYhglJsglContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/yhgl/jsgl/jsgl.js'
									]
								);
							}
						]
					}
				})
				
				/*   学生管理       */
				.state('app.jcsj.yhgl.xsgl', {
					url: '/xsgl',
					templateUrl: '../project/jcsj/html/jcsj/yhgl/xsgl/index.html',
					controller: 'jcsjYhglXsglContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/jcsj/css/jcsj/jcsj.css',
										'../project/jcsj/js/jcsj/yhgl/xsgl/xsgl.js'
									]
								);
							}
						]
					}
				})
				
				/*  //////////////////  基础数据 End  ///////////////////////  */




                    /*  /////////////////  中控面板   /////////////////////  */
				/*   中控面板     */
						.state('app.zkmb', {
							url: '/zkmb',
							templateUrl: '../project/zkmb/html/zkmb/index.html',
							controller: 'zkmbindexContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load([]).then(
												function() {
													return $ocLazyLoad.load(['../project/zkmb/css/zkmb/zkmb_v2.css',
														'../project/zkmb/js/zkmb/zkmb.js',
														'../vendor/jquery/blockUI/jquery.blockUI.js',
														'../vendor/echarts/echarts.min.js']);
												}
										);
									}
								]
							}
						})

						/*     出勤纠正       */
						.state('app.zkmb.kqjz', {
							url: '/kqjz',
							templateUrl: '../project/zkmb/html/zkmb/class/kqjz.html',
							controller: 'zkmbKqjzindexContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												['../project/zkmb/js/zkmb/zkmb.js']
										);
									}
								]
							}
						})


						/*     出勤统计        */
						.state('app.zkmb.kqtj', {
							url: '/kqtj',
							templateUrl: '../project/zkmb/html/zkmb/class/kqtj.html',
							controller: 'zkmbKqtjindexContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												['../project/zkmb/js/zkmb/zkmb.js']
										);
									}
								]
							}
						})


						/*     设备控制         */
						.state('app.zkmb.sbkz', {
							url: '/sbkz',
							templateUrl: '../project/zkmb/html/zkmb/sbkz/device.html',
							controller: 'zkmbSbkzindexContr',
							params:{"DeviceId":null,"page":null},
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												['../project/zkmb/js/zkmb/zkmb.js']
										);
									}
								]
							}
						})


			    /*  /////////////////  中控面板 End  /////////////////////  */


                /*  ////////////////////   权限管理   ////////////////////////   */
				//   权限管理-角色管理
				.state('app.qxgl', {
				    url: '/qxgl',
				    templateUrl: '../project/qxgl/html/qxgl/index.html',
				    controller: 'qxglContr',
				    resolve: {
				        deps: ['$ocLazyLoad',
							function ($ocLazyLoad) {
							    return $ocLazyLoad.load(
									['../project/qxgl/css/qxgl/qxgl.css', '../project/qxgl/js/qxgl/qxgl.js']
								);
							}
				        ]
				    }
				})

				//   权限管理-角色管理
				.state('app.qxgl.jsgl', {
				    url: '/jsgl',
				    templateUrl: '../project/qxgl/html/qxgl/jsgl/index.html',
				    controller: 'qxglJsglContr',
				    resolve: {
				        deps: ['$ocLazyLoad',
                            function ($ocLazyLoad) {
                                return $ocLazyLoad.load(
                                    ['../project/qxgl/js/qxgl/jsgl/jsgl.js']
                                );
                            }
				        ]
				    }
				})


				//   权限管理-模块管理
				.state('app.qxgl.mkgl', {
				    url: '/mkgl',
				    templateUrl: '../project/qxgl/html/qxgl/mkgl/index.html',
				    controller: 'qxglMkglContr',
				    resolve: {
				        deps: ['$ocLazyLoad',
                            function ($ocLazyLoad) {
                                return $ocLazyLoad.load(
                                    ['../project/qxgl/js/qxgl/mkgl/mkgl.js']
                                );
                            }
				        ]
				    }
				})


				//   权限管理-功能管理
				.state('app.qxgl.gngl', {
				    url: '/gngl',
				    templateUrl: '../project/qxgl/html/qxgl/gngl/index.html',
				    controller: 'qxglGnglContr',
				    resolve: {
				        deps: ['$ocLazyLoad',
                            function ($ocLazyLoad) {
                                return $ocLazyLoad.load(
                                    ['../project/qxgl/js/qxgl/gngl/gngl.js']
                                );
                            }
				        ]
				    }
				})


                //   权限管理-用户权限分配
				.state('app.qxgl.qxfp', {
				    url: '/qxfp',
				    templateUrl: '../project/qxgl/html/qxgl/qxfp/index.html',
				    controller: 'qxglQxfpContr',
				    resolve: {
				        deps: ['$ocLazyLoad',
							function ($ocLazyLoad) {
							    return $ocLazyLoad.load(
									['../project/qxgl/js/qxgl/qxfp/qxfp.js']
								);
							}
				        ]
				    }
				})


			    /*  ////////////////////   权限管理   End   ////////////////////////   */
			   
			   
			   
			   
			   
			   /*  //////////////////  出勤管理  ///////////////////////  */
				/*   出勤管理    */
				.state('app.cqgl', {
					url: '/cqgl',
					templateUrl: '../project/cqgl/html/cqgl/index.html',
					controller: 'cqglContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/cqgl/css/cqgl/cqgl.css',
										'../vendor/echarts/echarts.min.js',
										'../project/cqgl/js/cqgl/cqgl.js'
										
									]
								);
							}
						]
					}
				})
				
				/*   出勤统计    */
				.state('app.cqgl.cqtj', {
					url: '/cqtj',
					templateUrl: '../project/cqgl/html/cqgl/cqtj/index.html',
					controller: 'cqglCqtjContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/cqgl/css/cqgl/cqgl.css',
										'../vendor/echarts/echarts.min.js',
										'../project/cqgl/js/cqgl/cqtj/cqtj.js'
										
									]
								);
							}
						]
					}
				})
				
				/*   实时出勤    */
				.state('app.cqgl.sscq', {
					url: '/sscq',
					templateUrl: '../project/cqgl/html/cqgl/sscq/index.html',
					controller: 'cqglSscqContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/cqgl/css/cqgl/cqgl.css',
										'../project/cqgl/js/cqgl/sscq/sscq.js'
									]
								);
							}
						]
					}
				})
				/*   课堂出勤查看     */
				.state('app.cqgl.sscq.details', {
					url: '/details',
					templateUrl: '../project/cqgl/html/cqgl/sscq/details.html',
					controller: 'cqglSscqDetailsContr',
					params:{ClassroomId:null},
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/cqgl/css/cqgl/cqgl.css',
										'../project/cqgl/js/cqgl/sscq/sscq.js'
									]
								);
							}
						]
					}
				})
				/*   时时出勤-课堂出勤查看-查找列表     */
				.state('app.cqgl.sscq.details_list', {
					url: '/details_list',
					templateUrl: '../project/cqgl/html/cqgl/sscq/details_list.html',
					controller: 'cqglSscqDetailsListContr',
					params:{'ClassroomId':null,'be':"","end":""},
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/cqgl/css/cqgl/cqgl.css',
										'../project/cqgl/js/cqgl/sscq/sscq.js'
									]
								);
							}
						]
					}
				})
				
				/*   历史出勤    */
				.state('app.cqgl.lscq', {
					url: '/lscq',
					templateUrl: '../project/cqgl/html/cqgl/lscq/index.html',
					controller: 'cqglLscqContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/cqgl/css/cqgl/cqgl.css',
										'../project/cqgl/js/cqgl/lscq/lscq.js'
									]
								);
							}
						]
					}
				})
				
				/*   历史出勤 - 查看    */
				.state('app.cqgl.lscq.details', {
					url: '/details',
					templateUrl: '../project/cqgl/html/cqgl/lscq/details.html',
					controller: 'cqglLscqDetailsContr',
					params:{CId:null},
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/cqgl/css/cqgl/cqgl.css',
										'../project/cqgl/js/cqgl/lscq/lscq.js'
									]
								);
							}
						]
					}
				})
				/*  //////////////////  出勤管理 End  ///////////////////////  */
				
				
				
				
				/*  //////////////////  教室导流  ///////////////////////  */
				/*   教室导流      */
				.state('app.jsdl', {
					url: '/jsdl',
					templateUrl: '../project/jsdl/html/jsdl/index.html',
					controller: 'jsdlContr',
					resolve: {
							deps: ['$ocLazyLoad',
								function($ocLazyLoad) {
									return $ocLazyLoad.load(
										[
											'../project/jsdl/css/jsdl/jsdl.css',
											'../project/jsdl/js/jsdl/jsdl.js'
										]
									);
								}
							]
						}
				})
				
				/*   教室监控    */
				.state('app.jsdl.jsjk', {
					url: '/jsjk',
					templateUrl: '../project/jsdl/html/jsdl/jsjk/index.html',
					controller: 'jsdlJsjkContr',
					resolve: {
							deps: ['$ocLazyLoad',
								function($ocLazyLoad) {
									return $ocLazyLoad.load(
										[
											'../project/jsdl/css/jsdl/jsdl.css',
											'../project/jsdl/js/jsdl/jsjk/jsjk.js'
										]
									);
								}
							]
						}
				})
				
				
				/*   实时统计      */
				.state('app.jsdl.jsjk.sstj', {
					url: '/sstj',
					templateUrl: '../project/jsdl/html/jsdl/jsjk/sstj/index.html',
					controller: 'jsdlJsjkSstjContr',
					params:{ClassroomId:null},
					resolve: {
							deps: ['$ocLazyLoad',
								function($ocLazyLoad) {
									return $ocLazyLoad.load(
										[
											'../vendor/echarts/echarts.min.js',
											'../project/jsdl/css/jsdl/jsdl.css',
											'../project/jsdl/js/jsdl/jsjk/sstj/sstj.js'
										]
									);
								}
							]
						}
				})
				
				
				/*   导流分析    */
				.state('app.jsdl.dlfx', {
					url: '/dlfx',
					templateUrl: '../project/jsdl/html/jsdl/dlfx/index.html',
					controller: 'jsdlDlfxContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/jsdl/css/jsdl/jsdl.css',
										'../project/jsdl/js/jsdl/dlfx/rlt.js',
										'../project/jsdl/js/jsdl/dlfx/dlfx.js'
										
									]
								);
							}
						]
					}
				})
				/*  //////////////////  教室导流 End  ///////////////////////  */
				
				
				
				
				/*  //////////////////  课程管理   ///////////////////////  */
				/*   课程管理    */
				.state('app.kcgl', {
					url: '/kcgl',
					templateUrl: '../project/kcgl/html/kcgl/index.html',
					controller: 'kcglContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/kcgl/css/kcgl/kcgl.css',
										'../project/kcgl/js/kcgl/kcgl.js'
									]
								);
							}
						]
					}
				})
				
				
				/*   课程管理--查课表      */
				.state('app.kcgl.ckb', {
					url: '/ckb',
					templateUrl: '../project/kcgl/html/kcgl/ckb/index.html',
					controller: 'kcglCkbContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(['ui.select']).then(
									function() {
										return $ocLazyLoad.load(
											[
											'../project/kcgl/css/kcgl/kcgl.css',
											'../project/kcgl/js/kcgl/ckb/ckb.js'
											]
										);
									}
								);
							}
						]
					}
				})
				
				/*   课程管理--课程计划      */
				.state('app.kcgl.kcjh', {
					url: '/kcjh',
					templateUrl: '../project/kcgl/html/kcgl/kcjh/index.html',
					controller: 'kcglKcjhContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(['ui.select']).then(
									function() {
										return $ocLazyLoad.load(
											[
											'../project/kcgl/css/kcgl/kcgl.css',
											'../project/kcgl/js/kcgl/kcjh/kcjh.js'
											]
										);
									}
								);
							}
						]
					}
				})
                /*   课程管理--课程计划-章节配置列表      */
				.state('app.kcgl.kcjh.zjpz', {
				    url: '/zjpz',
				    templateUrl: '../project/kcgl/html/kcgl/kcjh/zjpz/index.html',
				    controller: 'kcglKcjhZjpzContr',
				    //    课程id , 课程名称 , 课程中间ID
				    params:{ccid:null,ccname:"",cctid:null},
				    resolve: {
				        deps: ['$ocLazyLoad',
							function ($ocLazyLoad) {
							    return $ocLazyLoad.load(['ui.select']).then(
									function () {
									    return $ocLazyLoad.load(
											[
											'../project/kcgl/css/kcgl/kcgl.css',
											'../project/kcgl/js/kcgl/kcjh/zjpz/zjpz.js'
											]
										);
									}
								);
							}
				        ]
				    }
				})


                /*   课程管理--基础课程      */
				.state('app.kcgl.jckc', {
				    url: '/jckc',
				    templateUrl: '../project/kcgl/html/kcgl/jckc/index.html',
				    controller: 'kcglJckcContr',
				    resolve: {
				        deps: ['$ocLazyLoad',
							function ($ocLazyLoad) {
							    return $ocLazyLoad.load(['ui.select']).then(
									function () {
									    return $ocLazyLoad.load(
											[
											'../project/kcgl/css/kcgl/kcgl.css',
											'../project/kcgl/js/kcgl/jckc/jckc.js'
											]
										);
									}
								);
							}
				        ]
				    }
				})

                /*   课程管理--基础课程 -- 添加      */
				.state('app.kcgl.jckc.add', {
				    url: '/add',
				    templateUrl: '../project/kcgl/html/kcgl/jckc/add/index.html',
				    controller: 'modalJckcAddContr',
				    params:{op:null,cid:null},
				    resolve: {
				        deps: ['$ocLazyLoad',
							function ($ocLazyLoad) {
							    return $ocLazyLoad.load(['ui.select']).then(
									function () {
									    return $ocLazyLoad.load(
											[
											'../project/kcgl/css/kcgl/kcgl.css',
											'../project/kcgl/js/kcgl/jckc/jckc.js'
											]
										);
									}
								);
							}
				        ]
				    }
				})
				
				/*   课程管理--课程计划--添加       */
				.state('app.kcgl.kcjh.add', {
					url: '/add',
					templateUrl: '../project/kcgl/html/kcgl/kcjh/add/index.html',
					controller: 'kcglKcjhAddContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
										'../project/kcgl/css/kcgl/kcgl.css',
										'../project/kcgl/js/kcgl/kcjh/add/add.js'
									]
								);
							}
						]
					}
				})
				/*  //////////////////  课程管理 End  ///////////////////////  */
				
				
				
				/*    ////////////   视频录播   ////////////    */
				//  视频录播
				.state('app.splb', {
					url: '/splb',
					templateUrl: '../project/splb/html/splb/index.html',
					controller: 'splbindexContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
									'../project/splb/css/splb/splb.css',
									'../project/splb/js/splb/splb.js'
									]
								);
							}
						]
					}
				})
				//  视频录播-详细页面
				.state('app.splb.xxym', {
					url: '/xxym',
					templateUrl: '../project/splb/html/splb/xxym/index.html',
					controller: 'splbXxymindexContr',
					params: {vid:"",cid:""},
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									[
									'../project/splb/css/splb/splb.css',
									'../project/splb/js/splb/xxym/xxym.js'
									]
								);
							}
						]
					}
				})
				//   视频录播 查看-课程详情
				.state('app.splbDetails', {
					url: '/splb/details',
					templateUrl: '../project/splb/html/splb/details.html',
					params: {
						Curriculumsid: null
					},
					controller: 'splbDetailsindexContr',
					resolve: {
						deps: ['$ocLazyLoad',
							function($ocLazyLoad) {
								return $ocLazyLoad.load(
									['../project/splb/js/splb/splb.js']
								);
							}
						]
					}
				})

				//   视频录播 查看-课程详情
				.state('app.detailsVideo', {
						url: '/splb/video',
						templateUrl: '../project/splb/html/splb/details_video.html',
						params: {
							Curriculumsid: null
						},
						controller: 'splbDetailsVideoindexContr',
						resolve: {
							deps: ['$ocLazyLoad',
								function($ocLazyLoad) {
									return $ocLazyLoad.load(
										['../project/splb/js/splb/splb.js']
									);
								}
							]
						}
					})
					/*    ////////////   视频录播   End     ////////////    */





			}
		]
	);