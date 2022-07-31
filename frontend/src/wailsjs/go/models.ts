export namespace bridge {
	
	export class DatabaseRequest {
	    conid: string;
	    database: string;
	
	    static createFrom(source: any = {}) {
	        return new DatabaseRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.conid = source["conid"];
	        this.database = source["database"];
	    }
	}
	export class DatabaseKeepOpenRequest {
	    conid: string;
	    database: string;
	    keepOpen: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DatabaseKeepOpenRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.conid = source["conid"];
	        this.database = source["database"];
	        this.keepOpen = source["keepOpen"];
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

export namespace serializer {
	
	export class Response {
	    status: number;
	    result: any;
	    message: string;
	    type: string;
	    time: number;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.result = source["result"];
	        this.message = source["message"];
	        this.type = source["type"];
	        this.time = source["time"];
	    }
	}

}

