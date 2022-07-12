export namespace bridge {
	
	export class DatabasePingRequest {
	    conid: string;
	    database: string;
	
	    static createFrom(source: any = {}) {
	        return new DatabasePingRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.conid = source["conid"];
	        this.database = source["database"];
	    }
	}
	export class ServerRefreshRequest {
	    conid: string;
	    keepOpen: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ServerRefreshRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.conid = source["conid"];
	        this.keepOpen = source["keepOpen"];
	    }
	}

}

