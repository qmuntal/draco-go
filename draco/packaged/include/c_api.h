// Copyright 2021 Quim Muntal.
//
// Licensed under the BSD 2-Clause License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://opensource.org/licenses/BSD-2-Clause
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
#ifndef DRACO_C_API_H_
#define DRACO_C_API_H_

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifdef _WIN32
# if defined(DRACO_C_BUILDING_DLL)
#  define EXPORT_API __declspec(dllexport)
# elif !defined(DRACO_C_STATIC)
#  define EXPORT_API __declspec(dllimport)
# else
#  define EXPORT_API
# endif
#elif __GNUC__ >= 4 || defined(__clang__)
# define EXPORT_API __attribute__((visibility ("default")))
#else
 #define EXPORT_API
#endif  // defined(_WIN32)

// draco::GeometryAttribute::Type

typedef enum {
    GT_INVALID = -1,
    GT_POSITION,
    GT_NORMAL,
    GT_COLOR,
    GT_TEX_COORD,
    GT_GENERIC
} dracoGeometryType;

// draco::DataType

typedef enum {
  DT_INVALID,
  DT_INT8,
  DT_UINT8,
  DT_INT16,
  DT_UINT16,
  DT_INT32,
  DT_UINT32,
  DT_INT64,
  DT_UINT64,
  DT_FLOAT32,
  DT_FLOAT64,
  DT_BOOL
} dracoDataType;

typedef const char* draco_string; // NULL terminated  

// draco::Status

typedef struct draco_status draco_status;

EXPORT_API void dracoStatusRelease(draco_status *status);

EXPORT_API int dracoStatusCode(const draco_status *status);

EXPORT_API bool dracoStatusOk(const draco_status *status);

EXPORT_API size_t dracoStatusErrorMsgLength(const draco_status *status);

EXPORT_API size_t dracoStatusErrorMsg(const draco_status *status, char *msg, size_t length);

// draco::PointAttribute

typedef struct draco_point_attr draco_point_attr;

EXPORT_API size_t dracoPointAttrSize(const draco_point_attr* pa);

EXPORT_API dracoGeometryType dracoPointAttrType(const draco_point_attr* pa);

EXPORT_API dracoDataType dracoPointAttrDataType(const draco_point_attr* pa);

EXPORT_API int8_t dracoPointAttrNumComponents(const draco_point_attr* pa);

EXPORT_API bool dracoPointAttrNormalized(const draco_point_attr* pa);

EXPORT_API int64_t dracoPointAttrByteStride(const draco_point_attr* pa);

EXPORT_API int64_t dracoPointAttrByteOffset(const draco_point_attr* pa);

EXPORT_API uint32_t dracoPointAttrUniqueId(const draco_point_attr* pa);

// draco::Mesh

typedef struct draco_mesh draco_mesh;

EXPORT_API draco_mesh* dracoNewMesh();

EXPORT_API void dracoMeshRelease(draco_mesh *mesh);

EXPORT_API uint32_t dracoMeshNumFaces(const draco_mesh *mesh);

EXPORT_API uint32_t dracoMeshNumPoints(const draco_mesh *mesh);

EXPORT_API int32_t dracoMeshNumAttrs(const draco_mesh *mesh);

// Queries an array of 3*face_count elements containing the triangle indices.
// out_values must be allocated to contain at least 3*face_count uint16_t elements.
// out_size must be exactly 3*face_count*sizeof(uint16_t), else out_values
// won´t be filled and returns false.
EXPORT_API bool dracoMeshGetTrianglesUint16(const draco_mesh *mesh,
                                            const size_t out_size,
                                            uint16_t *out_values);

// Queries an array of 3*face_count elements containing the triangle indices.
// out_values must be allocated to contain at least 3*face_count uint32_t elements.
// out_size must be exactly 3*face_count*sizeof(uint32_t), else out_values
// won´t be filled and returns false.
EXPORT_API bool dracoMeshGetTrianglesUint32(const draco_mesh *mesh,
                                            const size_t out_size,
                                            uint32_t *out_values);

EXPORT_API const draco_point_attr* dracoMeshGetAttribute(const draco_mesh *mesh, int32_t att_id);

EXPORT_API int32_t dracoMeshGetNamedAttributeId(const draco_mesh *mesh, dracoDataType data_type);

EXPORT_API const draco_point_attr* dracoMeshGetAttributeByUniqueId(const draco_mesh *mesh, uint32_t unique_id);

EXPORT_API bool dracoMeshGetAttributeData(const draco_mesh *mesh,
                                          const draco_point_attr *pa,
                                          dracoDataType data_type,
                                          const size_t out_size,
                                          void *out_values);

// draco::Decoder

typedef struct draco_decoder draco_decoder;

EXPORT_API draco_decoder* dracoNewDecoder();

EXPORT_API void dracoDecoderRelease(draco_decoder *decoder);

EXPORT_API draco_status* dracoDecoderArrayToMesh(draco_decoder *decoder, 
                                                 const char *data, 
                                                 size_t data_size,
                                                 draco_mesh *out_mesh);

#ifdef __cplusplus
}
#endif

#endif  // DRACO_C_API_H_