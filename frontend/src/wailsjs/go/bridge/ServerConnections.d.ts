// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {<-chan modules} from '../models';
import {bridge} from '../models';

export function Close(arg1:string,arg2:boolean):void;

export function ListDatabases(arg1:any):Promise<any>;

export function Listener(arg1:string,arg2:<-chan modules.EchoMessage):void;

export function Ping(arg1:bridge.PingRequest):Promise<any>;

export function Refresh(arg1:string):Promise<any>;

export function ServerStatus():Promise<any>;
