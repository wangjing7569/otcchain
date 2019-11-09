'use strict';

/**
 * Config for the router
 */
angular.module('app')
  .run(
    [          '$rootScope', '$state', '$stateParams',
      function ($rootScope,   $state,   $stateParams) {
          $rootScope.$state = $state;
          $rootScope.$stateParams = $stateParams;        
      }
    ]
  )
  .config(
    [          '$stateProvider', '$urlRouterProvider',
      function ($stateProvider,   $urlRouterProvider) {
          
          $urlRouterProvider
              .otherwise('/access/signin');
          $stateProvider
              .state('app', {
                  abstract: true,
                  url: '/app',
                  templateUrl: 'tpl/app.html'
              })
              .state('app.dashboard-v1', {
                  url: '/dashboard-v1',
                  templateUrl: 'tpl/app_dashboard_v1.html',
                  resolve: {
                    deps: ['$ocLazyLoad',
                      function( $ocLazyLoad ){
                        return $ocLazyLoad.load(['js/controllers/chart.js']);
                    }]
                  }
              })
              .state('app.dashboard-v2', {
                  url: '/dashboard-v2',
                  templateUrl: 'tpl/app_dashboard_v2.html',
                  resolve: {
                    deps: ['$ocLazyLoad',
                      function( $ocLazyLoad ){
                        return $ocLazyLoad.load(['js/controllers/chart.js']);
                    }]
                  }
              })

              // table
              .state('app.table', {
                  url: '/table',
                  template: '<div ui-view></div>'
              })

              .state('app.table.publishBill', {
                  url: '/publishBill',
                  templateUrl: 'tpl/page_publishBill.html'
              })

              .state('app.table.index',{
                  url: '/index',
                  templateUrl: 'tpl/page_index.html'
              })

              .state('app.table.myask',{
                  url: '/myask',
                  templateUrl: 'tpl/page_myask.html'
              })

              .state('app.table.myProvide',{
                  url: '/myProvide',
                  templateUrl: 'tpl/page_myProvide.html'
              })

              .state('app.table.transaction',{
                  url: '/trans',
                  templateUrl: 'tpl/page_transaction.html'
              })

              .state('app.table.calconstract',{
                  url: '/cal',
                  templateUrl: 'tpl/page_calconstract.html'
              })

              .state('app.table.hangconstract',{
                  url: '/hang',
                  templateUrl: 'tpl/page_hangconstract.html'
              })

              .state('app.table.historyconstract',{
                  url: '/hiscontract',
                  templateUrl: 'tpl/page_historyconstract.html'
              })


              .state('app.table.provideprice',{
                  url: '/provideprice',
                  templateUrl: 'tpl/provideprice.html'
              })


              .state('app.table.block',{
                  url: '/block',
                  templateUrl: 'tpl/page_blockExp.html'
              })

              .state('app.table.account',{
                  url: '/account',
                  templateUrl: 'tpl/page_accountManage.html'
              })

              .state('app.table.credit',{
                  url: '/credit',
                  templateUrl: 'tpl/page_creditManage.html'
              })


              .state('app.table.askprice',{
                  url: '/askprice',
                  templateUrl: 'tpl/askprice.html'
              })

              .state('app.table.myBill', {
                  url: '/myBill',
                  templateUrl: 'tpl/page_myBill.html'
              })

              .state('app.table.myUnBill', {
                  url: '/myUnBill',
                  templateUrl: 'tpl/page_myUnBill.html'
              })

              .state('app.table.BlockHistory', {
                  url: '/publishBill',
                  templateUrl: 'tpl/page_BlockHistory.html'
              })

              .state('app.page', {
                  url: '/page',
                  template: '<div ui-view></div>'
              })


              //access
              .state('access', {
                  url: '/access',
                  template: '<div ui-view class="fade-in-right-big smooth"></div>'
              })
              .state('access.signin', {
                  url: '/signin',
                  templateUrl: 'tpl/page_signin.html'
                  //,
                  //resolve: {
                  //    deps: ['uiLoad',
                  //      function( uiLoad ){
                  //        return uiLoad.load( ['js/controllers/signin.js'] );
                  //    }]
                  //}
              })
              .state('access.signup1', {
                  url: '/signup1',
                  templateUrl: 'tpl/page_signup1.html'//,
                  //resolve: {
                  //    deps: ['uiLoad',
                  //      function( uiLoad ){
                  //        return uiLoad.load( ['js/controllers/signup.js'] );
                  //    }]
                  //}
              })

              .state('access.signup2', {
                  url: '/signup2',
                  templateUrl: 'tpl/page_signup2.html'//,
                  //resolve: {
                  //    deps: ['uiLoad',
                  //      function( uiLoad ){
                  //        return uiLoad.load( ['js/controllers/signup.js'] );
                  //    }]
                  //}
              })

              .state('access.forgotpwd', {
                  url: '/forgotpwd',
                  templateUrl: 'tpl/page_forgotpwd.html'//,
                  //resolve: {
                  //    deps: ['uiLoad',
                  //      function( uiLoad ){
                  //        return uiLoad.load( ['js/controllers/signup.js'] );
                  //    }]
                  //}
              })

              .state('access.signup3', {
                  url: '/signup3',
                  templateUrl: 'tpl/page_signup3.html'
              })

              .state('access.404', {
                  url: '/404',
                  templateUrl: 'tpl/page_404.html'
              });

      }
    ]
  );