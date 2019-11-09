'use strict';

angular.module('app').constant(
		'REST_URL', 
			{
			    // login
			    'login': '/login',
			    'logout': '/logout',

				// invoke
				'invoke': '/channels/mychannel/chaincodes/otcc/invoke',
                // query
				'query': '/channels/mychannel/chaincodes/otcc/query',

				//invoke jiaoyi
				'invoke1': '/channels/mychannel/chaincodes/commute/invoke',

				//
				'query1': '/channels/mychannel/chaincodes/commute/query',

			}
	);