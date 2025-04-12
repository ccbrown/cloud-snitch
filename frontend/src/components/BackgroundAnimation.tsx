'use client';

import * as THREE from 'three';
import { MathUtils } from 'three';
import { MeshLineGeometry, MeshLineMaterial } from 'meshline';
import { extend, Canvas, useFrame } from '@react-three/fiber';
import { easing } from 'maath';
import React, { useEffect, useMemo, useRef, useState } from 'react';

extend({ MeshLineGeometry, MeshLineMaterial });

interface LinesProps {
    count: number;
    radius: number;
}

const Lines = ({ count, radius }: LinesProps) => {
    const lines = useMemo(() => {
        const MAX_WIDTH = radius / 80;
        return Array.from({ length: count }, () => {
            const pos = new THREE.Vector3(
                MathUtils.randFloatSpread(radius),
                MathUtils.randFloatSpread(radius),
                MathUtils.randFloatSpread(radius),
            );
            const points = [pos.clone(), pos.add(new THREE.Vector3(0, MathUtils.randFloatSpread(radius), 0))];
            return {
                color: Math.random() < 0.2 ? '#e47e7c' : '#9a72d0',
                width: MAX_WIDTH * Math.random(),
                speed: Math.max(0.05, 0.4 * Math.random()),
                points,
            };
        });
    }, [count, radius]);
    return lines.map((props, index) => <Line key={index} {...props} />);
};

interface LineProps {
    points: THREE.Vector3[];
    width: number;
    color: string;
    speed: number;
}

const Line = ({ points, width, color, speed }: LineProps) => {
    const ref = useRef<THREE.Mesh>(null);
    useFrame((_state, delta) => {
        if (ref.current) {
            (ref.current.material as MeshLineMaterial).dashOffset -= (delta * speed) / 10;
        }
    });
    return (
        <mesh ref={ref}>
            <meshLineGeometry points={points} />
            <meshLineMaterial
                transparent
                lineWidth={width}
                color={color}
                depthWrite={false}
                dashArray={0.5}
                dashRatio={0.95}
                toneMapped={false}
            />
        </mesh>
    );
};

interface RigProps {
    radius: number;
}

const Rig = ({ radius }: RigProps) => {
    useFrame((state, dt) => {
        const MOVE_FACTOR = -0.05;
        easing.damp3(
            state.camera.position,
            [
                (Math.sin(state.pointer.x * MOVE_FACTOR) * radius) / 2,
                (Math.atan(state.pointer.y * MOVE_FACTOR) * radius) / 2,
                radius / 2,
            ],
            0.25,
            dt,
        );
        state.camera.lookAt(0, 0, 0);
    });
    return null;
};

export const BackgroundAnimation = () => {
    const [eventSource, setEventSource] = useState<HTMLElement | undefined>(undefined);

    useEffect(() => {
        setEventSource(document.body);
    }, []);

    const RADIUS = 20;
    return (
        <Canvas eventSource={eventSource}>
            <color attach="background" args={['#dfdfdf']} />
            <Lines count={100} radius={RADIUS} />
            <Rig radius={RADIUS} />
        </Canvas>
    );
};
