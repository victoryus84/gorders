allprojects {
    repositories {
        google()
        mavenCentral()
    }
}

val fixedBuildDir = file("D:/Orders/frontend/build")

rootProject.layout.buildDirectory.value(
    rootProject.layout.projectDirectory.dir(fixedBuildDir.absolutePath)
)

subprojects {
    val newSubprojectBuildDir: Directory = 
        rootProject.layout.buildDirectory.get().dir(project.name)
    project.layout.buildDirectory.value(newSubprojectBuildDir)
}

subprojects {
    project.evaluationDependsOn(":app")
}

tasks.register<Delete>("clean") {
    delete(rootProject.layout.buildDirectory)
}