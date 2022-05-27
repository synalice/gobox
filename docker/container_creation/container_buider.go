package container_creation

type ContainerBuilder struct {
	container *Container
}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{container: &Container{}}
}

func (b *ContainerBuilder) SetConfigs() *ContainerConfigBuilder {
	return &ContainerConfigBuilder{*b}
}

func (b *ContainerBuilder) SetVolumes() *ContainerVolumeBuilder {
	return &ContainerVolumeBuilder{*b}
}

func (b *ContainerBuilder) ReturnNewContainer() *Container {
	return b.container
}
