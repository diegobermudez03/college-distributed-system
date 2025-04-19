output "server_exe_obj_name"{
    value = "${google_storage_bucket_object.server_exe.bucket}/${google_storage_bucket_object.server_exe.name}"
}

output "faculty_exe_obj_name"{
    value =  "${google_storage_bucket_object.faculty_exe.bucket}/${google_storage_bucket_object.faculty_exe.name}"
}

output "program_exe_obj_name"{
    value =  "${google_storage_bucket_object.program_exe.bucket}/${google_storage_bucket_object.program_exe.name}"
}

output "program_script_case1"{
    value =  "${google_storage_bucket_object.program_case1.bucket}/${google_storage_bucket_object.program_case1.name}"
}

output "faculty_script_case1"{
    value =  "${google_storage_bucket_object.faculty_case1.bucket}/${google_storage_bucket_object.faculty_case1.name}"
}


output "program_script_case2"{
    value =  "${google_storage_bucket_object.program_case2.bucket}/${google_storage_bucket_object.program_case2.name}"
}

output "faculty_script_case2"{
    value =  "${google_storage_bucket_object.faculty_case2.bucket}/${google_storage_bucket_object.faculty_case2.name}"
}

output "docker_compose_name"{
    value =  "${google_storage_bucket_object.docker_compose.bucket}/${google_storage_bucket_object.docker_compose.name}"
}

output "dti_exec_name"{
    value = google_storage_bucket_object.server_exe.name
}

output "fac_exec_name"{
    value = google_storage_bucket_object.faculty_exe.name
}

output "program_exec_name"{
    value = google_storage_bucket_object.program_exe.name
}