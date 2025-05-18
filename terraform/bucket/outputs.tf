####outputs for server exes (reqrep, lb, and proxy)
output "req_rep_exe_obj_name"{
    value = "${google_storage_bucket_object.req_rep_exe.bucket}/${google_storage_bucket_object.req_rep_exe.name}"
}

output "lb_exe_obj_name"{
    value = "${google_storage_bucket_object.lb_exe.bucket}/${google_storage_bucket_object.lb_exe.name}"
}

output "proxy_obj_name"{
    value = "${google_storage_bucket_object.lb_exe.bucket}/${google_storage_bucket_object.lb_exe.name}"
}

##### faculty and program exes
output "faculty_exe_obj_name"{
    value =  "${google_storage_bucket_object.faculty_exe.bucket}/${google_storage_bucket_object.faculty_exe.name}"
}

output "program_exe_obj_name"{
    value =  "${google_storage_bucket_object.program_exe.bucket}/${google_storage_bucket_object.program_exe.name}"
}

### scripts
output "script_case1"{
    value =  "${google_storage_bucket_object.case1.bucket}/${google_storage_bucket_object.case1.name}"
}

output "script_case1"{
    value =  "${google_storage_bucket_object.case2.bucket}/${google_storage_bucket_object.case2.name}"
}


##docker ocmpose
output "docker_compose_name"{
    value =  "${google_storage_bucket_object.docker_compose.bucket}/${google_storage_bucket_object.docker_compose.name}"
}

##### executables names
output "req_rep_exec_name"{
    value = google_storage_bucket_object.req_rep_exe.name
}

output "lb_exec_name"{
    value = google_storage_bucket_object.lb_exe.name
}

output "proxy_exec_name"{
    value = google_storage_bucket_object.proxy_exe.name
}

output "fac_exec_name"{
    value = google_storage_bucket_object.faculty_exe.name
}

output "program_exec_name"{
    value = google_storage_bucket_object.program_exe.name
}