package lda

import (
	"fmt"
	"math"

	"github.com/lda_api/internal/app/entity"
	"github.com/lda_api/internal/app/repository"
	"github.com/sirupsen/logrus"
)

type LDA struct {
	db *repository.DataBase
	//ConvAllData entity.ConversionData
	Alpha float64   `json:"alpha"`
	Beta  float64   `json:"beta"`
	X     []float64 `json:"X"`
	Y     []float64 `json:"Y"`
}

func New(db *repository.DataBase) *LDA {
	return &LDA{
		db: db,
	}
}

func (lda *LDA) ConvertData(dataUsers []entity.UserData) entity.ConversionData {
	var dataConv entity.ConversionData
	for i := 0; i < len(dataUsers); i++ {
		diffCoef := (float64)(dataUsers[i].LoanAmount) / float64(dataUsers[i].IncomeAnnum)
		diffCoef = (float64(dataUsers[i].IncomeAnnum) / diffCoef) / 1000
		if dataUsers[i].SelfEmployed == " Yes" {
			diffCoef = diffCoef * 0.25
		} else {
			diffCoef = diffCoef * 0.1
		}
		dataConv.ImportancecCoefficient = append(dataConv.ImportancecCoefficient, math.Round(diffCoef*100)/100)
		dataConv.Rating = append(dataConv.Rating, float64(dataUsers[i].CibilScore))
	}
	return dataConv
}

func GetMean(data []float64) float64 {
	if data == nil {
		return 0.0
	}
	var sum float64 = 0
	for i := 0; i < len(data); i++ {
		sum = sum + data[i]
	}

	return sum / float64(len(data))
}

func GetDispersion(data []float64, mean float64) float64 {
	if data == nil {
		return 0.0
	}
	var sum float64 = 0.0
	for i := 0; i < len(data); i++ {
		var temp float64 = data[i] - mean
		sum = sum + math.Pow(temp, 2)
	}
	return sum / (float64(len(data) - 1))
}

func GetCovariation(dataX []float64, meanX float64, dataY []float64, meanY float64) float64 {
	if dataX == nil && dataY == nil {
		fmt.Println("Error")
		return 0.0
	}
	var sum float64 = 0.0
	for i := 0; i < len(dataX); i++ {
		sum = sum + (dataX[i]-meanX)*(dataY[i]-meanY)
	}
	return sum / (float64(len(dataY) - 1))
}

func GetCovariationMatrix(mean []float64, data entity.ConversionData) [][]float64 {
	var covariationMatrix [][]float64
	covariation := GetCovariation(data.ImportancecCoefficient, mean[1], data.Rating, mean[0])
	fmt.Println(covariation)
	for i := 0; i < len(mean); i++ {
		var temp []float64
		for j := 0; j < len(mean); j++ {
			if i == j {
				if i == 0 {
					temp = append(temp, GetDispersion(data.Rating, mean[i]))
				} else {
					temp = append(temp, GetDispersion(data.ImportancecCoefficient, mean[i]))
				}
			} else {
				temp = append(temp, covariation)
			}
		}
		covariationMatrix = append(covariationMatrix, temp)
	}
	return covariationMatrix
}

func PlusMatrix(leftOperand, rightOperand [][]float64) [][]float64 {
	for i := 0; i < len(leftOperand); i++ {
		for j := 0; j < len(leftOperand[i]); j++ {
			leftOperand[i][j] = leftOperand[i][j] + rightOperand[i][j]
		}
	}
	return leftOperand
}

func SubtractionMatrix(leftOperand, rightOperand [][]float64) [][]float64 {
	for i := 0; i < len(leftOperand); i++ {
		for j := 0; j < len(leftOperand[i]); j++ {
			leftOperand[i][j] = leftOperand[i][j] - rightOperand[i][j]
		}
	}
	return leftOperand
}
func ProdMatrix(leftOperand, rightOperand [][]float64) [][]float64 {
	var resultMatrix [][]float64
	var sum float64 = 0.0

	for i := 0; i < len(leftOperand); i++ {
		var temp []float64
		for j := 0; j < len(rightOperand[i]); j++ {
			sum = 0
			for k := 0; k < len(leftOperand); k++ {
				sum = sum + (leftOperand[i][k] * rightOperand[k][j])
			}
			temp = append(temp, sum)
		}
		resultMatrix = append(resultMatrix, temp)
	}
	return resultMatrix
}

func GetDet(matrix [][]float64) float64 {
	var prod float64 = 1
	var prod2 float64 = 1
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			if i == j {
				prod = prod * matrix[i][j]
			} else {
				prod2 = prod2 * matrix[i][j]
			}
		}
	}
	return (prod - prod2)
}

func CheckUnitMatrix(inverseMatrix, W [][]float64) bool {
	result := ProdMatrix(W, inverseMatrix)
	var count int = 0
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i]); j++ {
			if result[i][j] >= 0.9 && result[i][j] <= 1 && i == j {
				count++
			} else {
				if int(result[i][j]) == 0 && i != j {
					count++
				}
			}
		}
	}
	fmt.Println(result)
	fmt.Println(count)
	if count == (len(result) * len(result[0])) {
		return true
	} else {
		return false
	}
}

func InverseMatrix(W [][]float64) [][]float64 {
	var det float64 = GetDet(W)
	fmt.Printf("det %f", det)
	var updW [][]float64
	for i := 1; i <= len(W); i += 1 {
		var tempForAdd []float64
		for j := 1; j <= len(W[0]); j += 1 {
			temp := math.Pow((-1), float64(i+j)) * W[len(W)-i][len(W[0])-j]
			tempForAdd = append(tempForAdd, temp)
		}
		updW = append(updW, tempForAdd)
	}
	fmt.Println("updW is find")
	fmt.Println(updW)
	var result [][]float64
	for i := 0; i < len(W); i++ {
		var tempForAdd []float64
		for j := 0; j < len(W[i]); j++ {
			temp := (1 / det) * updW[i][j]
			tempForAdd = append(tempForAdd, temp)
		}
		result = append(result, tempForAdd)
	}
	fmt.Print("result inverse")
	fmt.Println(result)
	if CheckUnitMatrix(result, W) {
		return result
	} else {
		logrus.Fatal("Матрица не является обратной")
		return nil
	}

}

func GetW(covariationMatrixFirstClass, covariationMatrixSecondClass [][]float64) [][]float64 {
	if covariationMatrixFirstClass == nil && covariationMatrixSecondClass == nil {
		return nil
	}
	W := PlusMatrix(covariationMatrixFirstClass, covariationMatrixSecondClass)
	for i := 0; i < len(covariationMatrixFirstClass); i++ {
		for j := 0; j < len(covariationMatrixFirstClass[i]); j++ {
			W[i][j] = W[i][j] / 2
		}
	}
	return W
}

func GetB(covariationMatrixAllData, W [][]float64) [][]float64 {
	return SubtractionMatrix(covariationMatrixAllData, W)
}

func GetS(W, B [][]float64) [][]float64 {
	fmt.Print("start W ")
	fmt.Println(W)
	W = InverseMatrix(W)
	fmt.Println("W inverse is find")
	fmt.Println(W)
	return ProdMatrix(W, B)
}

func (lda *LDA) FitModel() error {
	// инициализация массивов средних значений для классов и всего набора данных
	var meanValueFirstClass []float64
	var meanValueSecondClass []float64
	var meanValueAllData []float64
	/// инициализация ковариационных матриц для классов и всего набора данных
	var covariationMatrixAllData [][]float64
	var covariationMatrixFirstClass [][]float64
	var covariationMatrixSecondClass [][]float64
	// инициализация ковариационных матриц S, W, B
	var S [][]float64
	var W [][]float64
	var B [][]float64

	// Получение всех необходимых значений для всего набора данных
	dataAllUsers, err := lda.db.Data().SelectAllData()
	if err != nil {
		return err
	}
	fmt.Println(dataAllUsers)
	//convAllData := lda.ConvertData(dataAllUsers)
	var convAllData entity.ConversionData
	convAllData.ImportancecCoefficient = []float64{101, 89, 57, 76, 52, 49, 90, 72, 82, 88, 14, 33, 22, 20, 49, 36, 25, 29, 44, 42}
	convAllData.Rating = []float64{87, 71, 66, 55, 61, 66, 89, 79, 77, 59, 24, 34, 26, 39, 49, 42, 21, 44, 32, 39}
	meanValueAllData = append(meanValueAllData, GetMean(convAllData.Rating))
	meanValueAllData = append(meanValueAllData, GetMean(convAllData.ImportancecCoefficient))

	covariationMatrixAllData = GetCovariationMatrix(meanValueAllData, convAllData)

	fmt.Println(convAllData)
	fmt.Println(meanValueAllData)
	fmt.Println(covariationMatrixAllData)

	// Получение данных для класса "Approved"
	firstClass, err := lda.db.Data().SelectingDataByClass(" Approved")
	if err != nil {
		return err
	}
	fmt.Println(firstClass)
	//dataFirstClass := lda.ConvertData(firstClass)
	var dataFirstClass entity.ConversionData
	dataFirstClass.ImportancecCoefficient = []float64{101, 89, 57, 76, 52, 49, 90, 72, 82, 88}
	dataFirstClass.Rating = []float64{87, 71, 66, 55, 61, 66, 89, 79, 77, 59}
	meanValueFirstClass = append(meanValueFirstClass, GetMean(dataFirstClass.Rating))
	meanValueFirstClass = append(meanValueFirstClass, GetMean(dataFirstClass.ImportancecCoefficient))

	covariationMatrixFirstClass = GetCovariationMatrix(meanValueFirstClass, dataFirstClass)

	fmt.Print("Данные класса одобрено: ")
	fmt.Println(dataFirstClass)
	fmt.Println(meanValueFirstClass)
	fmt.Println(covariationMatrixFirstClass)

	// Получение данных для класса "Rejected"
	secondClass, err := lda.db.Data().SelectingDataByClass(" Rejected")
	if err != nil {
		return err
	}
	fmt.Println(secondClass)
	//dataSecondClass := lda.ConvertData(secondClass)
	var dataSecondClass entity.ConversionData
	dataSecondClass.ImportancecCoefficient = []float64{14, 33, 22, 20, 49, 36, 25, 29, 44, 42}
	dataSecondClass.Rating = []float64{24, 34, 26, 39, 49, 42, 21, 44, 32, 39}
	meanValueSecondClass = append(meanValueSecondClass, GetMean(dataSecondClass.Rating))
	meanValueSecondClass = append(meanValueSecondClass, GetMean(dataSecondClass.ImportancecCoefficient))

	covariationMatrixSecondClass = GetCovariationMatrix(meanValueSecondClass, dataSecondClass)

	fmt.Print("Данные класса отклонено: ")
	fmt.Println(dataSecondClass)
	fmt.Println(meanValueSecondClass)
	fmt.Println(covariationMatrixSecondClass)

	// Произведение расчетов для ковариационных матриц S, B, W

	W = GetW(covariationMatrixFirstClass, covariationMatrixSecondClass)
	fmt.Println("W is find")
	B = GetB(covariationMatrixAllData, W)
	fmt.Println("B is fund")
	S = GetS(W, B)
	fmt.Println(S)
	return nil
}
