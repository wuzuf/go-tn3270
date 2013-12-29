package tn3270

type Parser interface {
    Parse([] byte) error
}

type parser struct {
	state state
	tnh TNHandler
	tn3270negoh TN3270NegoHandler
	tn3270h TN3270Handler
	errorh ErrorHandler
}

type state struct {
    data []byte
    length int
    position int
    stack [10]int
    count int
    top int
    ts int
    cs int
    name *[]byte

    starttxt int
    
    idx int
    addr [2]byte
    resourceName []byte
    deviceName []byte
    deviceType []byte
    functionsList []byte
    
}

func (parser *parser) StartTxt() {
    parser.state.starttxt = parser.state.position
}

func (parser *parser) EndTxt() {
    if (parser.state.starttxt != -1 && parser.state.starttxt < parser.state.position) {
        parser.tn3270h.OnTN3270Text(parser.state.data[parser.state.starttxt:parser.state.position])
    }
    parser.state.starttxt = -1
}

func (state *state) GetAddr() int {
    if ((state.addr[0] & 0xC0) == 0x00)  {
        return (int(state.addr[0] & 0x3F) << 8) | int(state.addr[1])
    }
    return (int(state.addr[0] & 0x3F) << 6) | int(state.addr[1] & 0x3F)
}

%%{
    machine tn3270;
    alphtype byte;
    access state.;
    variable p state.position;
    variable pe state.length;

    action error { parser.errorh.OnError(state.data, state.position); }

    action tn_command { parser.tnh.OnTNCommand(fc);}
    action tn_argcommand { parser.tnh.OnTNArgCommand(state.data[state.position-1], fc);}

    action tn3270_command { parser.tn3270h.OnTN3270Command(fc); }
    action tn3270_aid { parser.tn3270h.OnTN3270AID(fc); }
    action tn3270_wcc { parser.tn3270h.OnTN3270WCC(fc); }
    action tn3270_sba { parser.tn3270h.OnTN3270SBA(state.GetAddr()); }
    action tn3270_eua { parser.tn3270h.OnTN3270EUA(state.GetAddr()); }
    action tn3270_ic { parser.tn3270h.OnTN3270IC(); }
    action tn3270_pt { parser.tn3270h.OnTN3270PT(); }
    action tn3270_sf { parser.tn3270h.OnTN3270SF(fc); }
    action tn3270_ra { parser.tn3270h.OnTN3270RA(state.GetAddr(), fc); }
    action tn3270_sfe { parser.tn3270h.OnTN3270SFE(fc); state.count = int(fc); fcall tn3270_args; }
    action tn3270_message { parser.EndTxt(); parser.tn3270h.OnTN3270Message(); }

    action tn3270_resource_name {
        state.name = &state.resourceName
    }
    action tn3270_addr {
    	addr := state.addr[:][:0]
        state.name = &addr
    }
    action tn3270_device_name {
        state.name = &state.deviceName
    }
    action tn3270_device_type {
        state.name = &state.deviceType
    }
    action tn3270_functions_list {
        state.name = &state.functionsList
    }
    action tn3270_name {
    	*state.name = append(*state.name, fc)
    }
    
    action tn3270_name_end {
        state.name = nil
    }

    action tn3270_function_request {
        parser.tn3270negoh.OnTN3270FunctionsRequest(state.functionsList);
    }
    action tn3270_function_is {
        parser.tn3270negoh.OnTN3270FunctionsIs(state.functionsList);
    }
    action tn3270_send_device_type {
        parser.tn3270negoh.OnTN3270SendDeviceType();
    }
    action tn3270_device_type_request {
        parser.tn3270negoh.OnTN3270DeviceTypeRequest(state.deviceType, state.deviceName, state.resourceName);
    }
    action tn3270_device_type_is {
        parser.tn3270negoh.OnTN3270DeviceTypeIs(state.deviceType, state.deviceName);
    }
    action tn3270_device_type_reject {
        parser.tn3270negoh.OnTN3270DeviceTypeReject(fc);
    }
    action tn3270_header {}

    action tn3270_starttxt { parser.StartTxt(); }
    action tn3270_endtxt { parser.EndTxt(); }

    action tn3270_endarg { state.count--; if(state.count == 0) { fret; } }

    action tn_subneg { fcall tn3270_subneg; }
    action tn_subneg_end { fret; }

    ##########
    # TELNET
    ##########

    tn_iac = 0xff;
    tn_se = 240;
    tn_nop = 241;
    tn_dm = 242;
    tn_brk = 243;
    tn_ip = 244;
    tn_ao = 245;
    tn_ayt = 246;
    tn_ec = 247;
    tn_el = 248;
    tn_ga = 249;
    tn_sb = 250;
    tn_will = 251;
    tn_wont = 252;
    tn_do = 253;
    tn_dont = 254;
    tn_eor = 239;
      
    tn_command = tn_nop | tn_brk | tn_ip | tn_ao | tn_ayt | tn_ec | tn_el | tn_ga;
    tn_command_arg = tn_will | tn_wont | tn_do | tn_dont;
    tn_commmand_subneg = tn_sb;

    tn_plain_text = (^tn_iac | tn_iac tn_iac);
      
    tn_basic_command  = tn_iac tn_command @tn_command;
    tn_arg_command    = tn_iac tn_command_arg any @tn_argcommand ;
    tn_subneg_command = tn_iac tn_commmand_subneg any @tn_subneg;

    tn_iac_sequence = ( tn_basic_command | tn_arg_command | tn_subneg_command );

    ##########
    # TN3270E
    ##########

    tn3270_tn3270e             = 0x28;
    tn3270_associate           = 0x00;
    tn3270_connect             = 0x01;
    tn3270_device_type         = 0x02;
    tn3270_functions           = 0x03;
    tn3270_is                  = 0x04;
    tn3270_reason              = 0x05;
    tn3270_reject              = 0x06;
    tn3270_request             = 0x07;
    tn3270_send                = 0x08;
    tn3270_conn_partner        = 0x00;
    tn3270_device_in_use       = 0x01;
    tn3270_inv_associate       = 0x02;
    tn3270_inv_name            = 0x03;
    tn3270_inv_device_type     = 0x04;
    tn3270_type_name_error     = 0x05;
    tn3270_unknown_error       = 0x06;
    tn3270_unsupported_req     = 0x07;
    tn3270_bind_image          = 0x00;
    tn3270_data_stream_ctl     = 0x01;
    tn3270_responses           = 0x02;
    tn3270_scs_ctl_codes       = 0x03;
    tn3270_sysreq              = 0x04;

    tn3270_ascii = [A-Z0-9_\-];

    tn3270_resource_name = tn3270_ascii{1,8} >tn3270_resource_name $tn3270_name  %tn3270_name_end;
    tn3270_device_types = tn3270_ascii{1,15} >tn3270_device_type   $tn3270_name  %tn3270_name_end;
    tn3270_device_name = tn3270_ascii{1,8}   >tn3270_device_name   $tn3270_name  %tn3270_name_end;
    tn3270_reason_code = tn3270_conn_partner | tn3270_inv_associate | tn3270_inv_name | tn3270_inv_device_type | tn3270_type_name_error | tn3270_unknown_error | tn3270_unsupported_req;
    tn3270_function_list = ( tn3270_bind_image | tn3270_data_stream_ctl | tn3270_responses | tn3270_scs_ctl_codes | tn3270_sysreq ){0,20} >tn3270_functions_list  $tn3270_name  %tn3270_name_end;

    # subnegociation
    tn3270_subneg_send_device_type = tn3270_send . tn3270_device_type %tn3270_send_device_type;
    tn3270_subneg_device_type_request = tn3270_device_type . tn3270_request . tn3270_device_types . ( tn3270_connect . tn3270_resource_name | tn3270_associate . tn3270_device_name ) %tn3270_device_type_request;
    tn3270_subneg_device_type_is = tn3270_device_type . tn3270_is . tn3270_device_types . tn3270_connect . tn3270_device_name %tn3270_device_type_is;
    tn3270_subneg_device_type_reject = tn3270_device_type . tn3270_reject . tn3270_reason . tn3270_reason_code %tn3270_device_type_reject;
    tn3270_subneg_function_request = tn3270_functions . tn3270_request . tn3270_function_list %tn3270_function_request;
    tn3270_subneg_function_is = tn3270_functions . tn3270_is . tn3270_function_list %tn3270_function_is;
    tn3270_subneg_list = tn3270_subneg_send_device_type
                       | tn3270_subneg_device_type_request
                       | tn3270_subneg_device_type_is
                       | tn3270_subneg_device_type_reject
                       | tn3270_subneg_function_request
                       | tn3270_subneg_function_is;
    
    tn3270_subneg := tn3270_subneg_list . tn_iac . tn_se @tn_subneg_end;

    tn3270_arg = any.any @tn3270_endarg;
    tn3270_args := tn3270_arg+;

    tn3270_command = (0x05 | 0xf5 | 0x01 | 0xf1 | 0x7e | 0x6f | 0xf6 | 0x6e | 0xf2 | 0xf3) @tn3270_command;
    tn3270_wcc = any @tn3270_wcc;
    tn3270_aid = 0x7d @tn3270_aid;
    tn3270_addr = any{2} >tn3270_addr  $tn3270_name  %tn3270_name_end;

    # orders
    tn3270_sba = 0x11 . tn3270_addr @tn3270_sba;
    tn3270_sf = 0x1d . any @tn3270_sf;
    tn3270_ic = 0x13 @tn3270_ic;
    tn3270_eua = 0x12 . tn3270_addr @tn3270_eua;
    tn3270_pt = 0x05 @tn3270_pt;
    tn3270_sfe = 0x29 . any @tn3270_sfe;
    tn3270_ra = 0x3c . tn3270_addr . any @tn3270_ra;
        
    tn3270_order = ( tn3270_sba | tn3270_sf | tn3270_ic | tn3270_eua | tn3270_pt | tn3270_sfe | tn3270_ra ) >tn3270_endtxt %tn3270_starttxt;
    tn3270_plain_text = (any - (0x11 | 0x1d | 0x12 | 0x05 | 0x29 | 0x3c | tn_iac)) +;
    tn3270_content = (tn3270_order | tn3270_plain_text) *;
    tn3270_header = any {5} @tn3270_header;
    tn3270_data = ( ( (tn3270_command . tn3270_wcc) | (tn3270_aid . tn3270_addr) ) . tn3270_content);
    tn3270_message = tn3270_header . tn3270_data . tn_iac @tn3270_message;
    main := ( tn_iac_sequence | tn3270_message . tn_eor )*  $err(error);

}%%

%% write data;

func (parser *parser) Init() {
	state := &parser.state
    state.starttxt = -1
    
    %% write init;
}

func (parser *parser) Parse(data []byte ) error {
	state := &parser.state
    state.position = 0
    state.data = data
    state.length = len(data)
    eof := 0

    %% write exec;

    // Store any pending text
    parser.EndTxt()
    
    return nil
}

func (state *state) Finish () int {
    if ( state.cs == tn3270_error ) {
        return -1
    }
    if ( state.cs >= tn3270_first_final ) {
        return 1
    }
    return 0
}
