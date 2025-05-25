####outputs for server exes (reqrep, lb, and proxy)
output "req_rep_exe_obj_name"{
    value = "${google_storage_bucket_object.req_rep_exe.bucket}/${google_storage_bucket_object.req_rep_exe.name}"
}

output "lb_exe_obj_name"{
    value = "${google_storage_bucket_object.lb_exe.bucket}/${google_storage_bucket_object.lb_exe.name}"
}

output "proxy_obj_name"{
    value = "${google_storage_bucket_object.proxy_exe.bucket}/${google_storage_bucket_object.proxy_exe.name}"
}

##### faculty and program exes
output "faculty_exe_obj_name"{
    value =  "${google_storage_bucket_object.faculty_exe.bucket}/${google_storage_bucket_object.faculty_exe.name}"
}

output "program_exe_obj_name"{
    value =  "${google_storage_bucket_object.program_exe.bucket}/${google_storage_bucket_object.program_exe.name}"
}

###################################################################### scripts
output "scase1_fac"{
    value =  "${google_storage_bucket_object.case1fac.bucket}/${google_storage_bucket_object.case1fac.name}"
}
output "scase1_programs"{
    value =  "${google_storage_bucket_object.case1programs.bucket}/${google_storage_bucket_object.case1programs.name}"
}

output "scase2_fac"{
    value =  "${google_storage_bucket_object.case2fac.bucket}/${google_storage_bucket_object.case2fac.name}"
}
output "scase2_programs"{
    value =  "${google_storage_bucket_object.case2programs.bucket}/${google_storage_bucket_object.case2programs.name}"
}


output "scase3_fac"{
    value =  "${google_storage_bucket_object.case3fac.bucket}/${google_storage_bucket_object.case3fac.name}"
}
output "scase3_programs"{
    value =  "${google_storage_bucket_object.case3programs.bucket}/${google_storage_bucket_object.case3programs.name}"
}


output "scase4_fac"{
    value =  "${google_storage_bucket_object.case4fac.bucket}/${google_storage_bucket_object.case4fac.name}"
}
output "scase4_programs"{
    value =  "${google_storage_bucket_object.case4programs.bucket}/${google_storage_bucket_object.case4programs.name}"
}

output "scase5_fac"{
    value =  "${google_storage_bucket_object.case5fac.bucket}/${google_storage_bucket_object.case5fac.name}"
}
output "scase5_programs1"{
    value =  "${google_storage_bucket_object.case5programs1.bucket}/${google_storage_bucket_object.case5programs1.name}"
}
output "scase5_programs2"{
    value =  "${google_storage_bucket_object.case5programs2.bucket}/${google_storage_bucket_object.case5programs2.name}"
}

output "bash_fac"{
    value =  "${google_storage_bucket_object.bashfac.bucket}/${google_storage_bucket_object.bashfac.name}"
}
output "bash_programs"{
    value =  "${google_storage_bucket_object.bashprograms.bucket}/${google_storage_bucket_object.bashprograms.name}"
}



##################################################################################################################

output "metricasc1"{
    value =  "${google_storage_bucket_object.metricasc1.bucket}/${google_storage_bucket_object.metricasc1.name}"
}
output "metricasc2"{
    value =  "${google_storage_bucket_object.metricasc2.bucket}/${google_storage_bucket_object.metricasc2.name}"
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