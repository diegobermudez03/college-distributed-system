output "server_exe_obj_name"{
    value = google_storage_bucket_object.server_exe.name
}

output "faculty_exe_obj_name"{
    value = google_storage_bucket_object.faculty_exe.name
}

output "program_exe_obj_name"{
    value = google_storage_bucket_object.program_exe.name
}

output "program_script_name"{
    value = google_storage_bucket_object.program_script.name
}

output "faculty_script_name"{
    value = google_storage_bucket_object.faculty_script.name
}

output "docker_compose_name"{
    value = google_storage_bucket_object.docker_compose.name
}