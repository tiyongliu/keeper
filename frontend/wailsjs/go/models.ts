export namespace serializer {
	
	export class Response {
	    code: number;
	    result?: any;
	    message: string;
	    type: string;
	    time: number;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.result = source["result"];
	        this.message = source["message"];
	        this.type = source["type"];
	        this.time = source["time"];
	    }
	}

}

