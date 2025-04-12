import { MeshLineGeometry, MeshLineMaterial } from 'meshline';
import { extend, Object3DNode, MaterialNode } from '@react-three/fiber';

declare module '@react-three/fiber' {
    interface ThreeElements {
        meshLineGeometry: Object3DNode<MeshLineGeometry, typeof MeshLineGeometry>;
        meshLineMaterial: MaterialNode<MeshLineMaterial, typeof MeshLineMaterial>;
    }
}

extend({ MeshLineGeometry, MeshLineMaterial });
